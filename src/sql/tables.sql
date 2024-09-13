CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE moderation_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type service_type,
    moderation_status moderation_status,
    version_id int default 1,
    organization_responsible_id UUID REFERENCES organization_responsible(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tender_version (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    version_id INT NOT NULL,
    value jsonb NOT NULL
);


CREATE OR REPLACE FUNCTION log_tender_changes() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        INSERT INTO tender_version (tender_id, version_id, value)
        VALUES (
            OLD.id,  
            OLD.version_id + 1,
            to_jsonb(OLD)
        );
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO tender_version (tender_id, version_id, value)
        VALUES (
            NEW.id,  
            1,
            to_jsonb(NEW)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER tender_change_trigger
AFTER INSERT OR UPDATE ON tender
FOR EACH ROW EXECUTE FUNCTION log_tender_changes();

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

CREATE TYPE author_type AS ENUM (
    'Organization',
    'User'
);


CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    bid_status bid_status,
    version_id int default 1,
    author_type author_type not null
    author_id UUID REFERENCES organization_responsible(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tender_id UUID REFERENCES tender(id)
);

CREATE TABLE bid_version (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    version_id INT NOT NULL,
    value jsonb NOT NULL
);

CREATE OR REPLACE FUNCTION log_bid_changes() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        INSERT INTO bid_version (bid_id, version_id, value)
        VALUES (
            OLD.id,  
            OLD.version_id + 1,
            to_jsonb(OLD)
        );
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO bid_version (bid_id, version_id, value)
        VALUES (
            NEW.id,  
            1,
            to_jsonb(NEW)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TYPE decision AS ENUM (
    'Approved',
    'Rejected',
);

CREATE TABLE bid_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    value decision decision NOT NULL
);

CREATE TABLE bid_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    value TEXT NOT NULL
);