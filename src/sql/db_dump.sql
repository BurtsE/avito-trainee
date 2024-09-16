--
-- PostgreSQL database dump
--

-- Dumped from database version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.13 (Ubuntu 14.13-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: author_type; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.author_type AS ENUM (
    'Organization',
    'User'
);


ALTER TYPE public.author_type OWNER TO admin;

--
-- Name: bid_status; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);


ALTER TYPE public.bid_status OWNER TO admin;

--
-- Name: decision; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.decision AS ENUM (
    'Approved',
    'Rejected'
);


ALTER TYPE public.decision OWNER TO admin;

--
-- Name: moderation_status; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.moderation_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);


ALTER TYPE public.moderation_status OWNER TO admin;

--
-- Name: organization_type; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);


ALTER TYPE public.organization_type OWNER TO admin;

--
-- Name: service_type; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);


ALTER TYPE public.service_type OWNER TO admin;

--
-- Name: log_bid_changes(); Type: FUNCTION; Schema: public; Owner: admin
--

CREATE FUNCTION public.log_bid_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.log_bid_changes() OWNER TO admin;

--
-- Name: log_tender_changes(); Type: FUNCTION; Schema: public; Owner: admin
--

CREATE FUNCTION public.log_tender_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.log_tender_changes() OWNER TO admin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bid; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bid (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    bid_status public.bid_status DEFAULT 'Created'::public.bid_status,
    version_id integer DEFAULT 1,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    author_id uuid NOT NULL,
    tender_id uuid,
    author_type public.author_type NOT NULL
);


ALTER TABLE public.bid OWNER TO admin;

--
-- Name: bid_decisions; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bid_decisions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    bid_id uuid,
    value public.decision
);


ALTER TABLE public.bid_decisions OWNER TO admin;

--
-- Name: bid_feedback; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bid_feedback (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    bid_id uuid,
    value text NOT NULL
);


ALTER TABLE public.bid_feedback OWNER TO admin;

--
-- Name: bid_version; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bid_version (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    bid_id uuid,
    version_id integer NOT NULL,
    value jsonb NOT NULL
);


ALTER TABLE public.bid_version OWNER TO admin;

--
-- Name: employee; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.employee (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(50) NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.employee OWNER TO admin;

--
-- Name: organization; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.organization (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    type public.organization_type,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.organization OWNER TO admin;

--
-- Name: organization_responsible; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.organization_responsible (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    organization_id uuid,
    user_id uuid
);


ALTER TABLE public.organization_responsible OWNER TO admin;

--
-- Name: tender; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.tender (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    service_type public.service_type,
    version_id integer DEFAULT 1,
    organization_responsible_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    moderation_status public.moderation_status
);


ALTER TABLE public.tender OWNER TO admin;

--
-- Name: tender_version; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.tender_version (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    tender_id uuid,
    version_id integer NOT NULL,
    value jsonb NOT NULL
);


ALTER TABLE public.tender_version OWNER TO admin;

--
-- Data for Name: bid; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.bid (id, name, description, bid_status, version_id, created_at, updated_at, author_id, tender_id, author_type) FROM stdin;
48dc4e16-45b6-4741-bde9-05c3a65ec63e	tender1 bid 1	string	Created	1	2024-09-13 16:36:32.263791	2024-09-13 16:36:32.263791	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	8ad46627-f6bb-4bc8-968a-14c6a7a8b99b	Organization
f5e8acc8-c329-4309-b56a-44fddbb40b16	tender1 bid 2	string	Created	1	2024-09-13 16:36:36.760562	2024-09-13 16:36:36.760562	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	8ad46627-f6bb-4bc8-968a-14c6a7a8b99b	Organization
d0aff544-9998-48a4-8eef-6c4179686564	tender2 bid 1	string	Created	1	2024-09-13 16:37:10.069946	2024-09-13 16:37:10.069946	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	Organization
226fa3e9-4cf3-4d9a-830f-9d7163632ff9	tender2 bid 2	string	Created	1	2024-09-13 16:37:18.558527	2024-09-13 16:37:18.558527	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	Organization
1e982f10-9e4a-461a-be01-f0b9919247ff	tender2 bid 3	string	Created	1	2024-09-13 16:38:51.474458	2024-09-13 16:38:51.474458	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	Organization
7303dc6b-7145-485e-b1cb-13ec536215a7	tender2 bid 4	string	Created	1	2024-09-13 16:39:59.968456	2024-09-13 16:39:59.968456	a6cac94a-4023-4719-b7fe-5ef7205cacd7	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	User
ce4a07a1-1f2e-4bc6-a635-92090d17345f	tender2 bid 5	string	Created	4	2024-09-13 16:40:04.561062	2024-09-13 16:41:15.58714	a6cac94a-4023-4719-b7fe-5ef7205cacd7	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	User
\.


--
-- Data for Name: bid_decisions; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.bid_decisions (id, bid_id, value) FROM stdin;
e3b49d51-039e-47a1-be3a-7eb4af8470ba	f5e8acc8-c329-4309-b56a-44fddbb40b16	Approved
\.


--
-- Data for Name: bid_feedback; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.bid_feedback (id, bid_id, value) FROM stdin;
dbd1b335-fd97-4195-b0e9-e94c07dd04b3	ce4a07a1-1f2e-4bc6-a635-92090d17345f	
fb6aae79-ef6d-4f5e-bdc7-59bf037a1b91	ce4a07a1-1f2e-4bc6-a635-92090d17345f	
4aa5a2a2-36d5-44aa-918e-0ec256ae4c08	ce4a07a1-1f2e-4bc6-a635-92090d17345f	
45738d24-ade6-42e8-aba1-bd934f15a790	ce4a07a1-1f2e-4bc6-a635-92090d17345f	cool
\.


--
-- Data for Name: bid_version; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.bid_version (id, bid_id, version_id, value) FROM stdin;
922762d5-5a43-4054-9ca2-3af2230bfd10	48dc4e16-45b6-4741-bde9-05c3a65ec63e	1	{"id": "48dc4e16-45b6-4741-bde9-05c3a65ec63e", "name": "tender1 bid 1", "author_id": "37e1fd6c-9863-4052-8dd0-c4cc00bcb25d", "tender_id": "8ad46627-f6bb-4bc8-968a-14c6a7a8b99b", "bid_status": "Created", "created_at": "2024-09-13T16:36:32.263791", "updated_at": "2024-09-13T16:36:32.263791", "version_id": 1, "author_type": "Organization", "description": "string"}
55c03c7b-6669-4a6d-b46c-082ac6c3316f	f5e8acc8-c329-4309-b56a-44fddbb40b16	1	{"id": "f5e8acc8-c329-4309-b56a-44fddbb40b16", "name": "tender1 bid 2", "author_id": "37e1fd6c-9863-4052-8dd0-c4cc00bcb25d", "tender_id": "8ad46627-f6bb-4bc8-968a-14c6a7a8b99b", "bid_status": "Created", "created_at": "2024-09-13T16:36:36.760562", "updated_at": "2024-09-13T16:36:36.760562", "version_id": 1, "author_type": "Organization", "description": "string"}
0d160d49-1494-4015-bfa4-0febf42c7209	d0aff544-9998-48a4-8eef-6c4179686564	1	{"id": "d0aff544-9998-48a4-8eef-6c4179686564", "name": "tender2 bid 1", "author_id": "37e1fd6c-9863-4052-8dd0-c4cc00bcb25d", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:37:10.069946", "updated_at": "2024-09-13T16:37:10.069946", "version_id": 1, "author_type": "Organization", "description": "string"}
876e2115-171f-4e4c-80fb-8737d388d489	226fa3e9-4cf3-4d9a-830f-9d7163632ff9	1	{"id": "226fa3e9-4cf3-4d9a-830f-9d7163632ff9", "name": "tender2 bid 2", "author_id": "37e1fd6c-9863-4052-8dd0-c4cc00bcb25d", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:37:18.558527", "updated_at": "2024-09-13T16:37:18.558527", "version_id": 1, "author_type": "Organization", "description": "string"}
5d5d7f58-435d-434a-871a-2e1e18383220	1e982f10-9e4a-461a-be01-f0b9919247ff	1	{"id": "1e982f10-9e4a-461a-be01-f0b9919247ff", "name": "tender2 bid 3", "author_id": "37e1fd6c-9863-4052-8dd0-c4cc00bcb25d", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:38:51.474458", "updated_at": "2024-09-13T16:38:51.474458", "version_id": 1, "author_type": "Organization", "description": "string"}
31a6f303-e878-445c-8f7a-b05933be58e3	7303dc6b-7145-485e-b1cb-13ec536215a7	1	{"id": "7303dc6b-7145-485e-b1cb-13ec536215a7", "name": "tender2 bid 4", "author_id": "a6cac94a-4023-4719-b7fe-5ef7205cacd7", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:39:59.968456", "updated_at": "2024-09-13T16:39:59.968456", "version_id": 1, "author_type": "User", "description": "string"}
466bc89e-c66c-4066-8a18-fa40dd3ce0af	ce4a07a1-1f2e-4bc6-a635-92090d17345f	1	{"id": "ce4a07a1-1f2e-4bc6-a635-92090d17345f", "name": "tender2 bid 5", "author_id": "a6cac94a-4023-4719-b7fe-5ef7205cacd7", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:40:04.561062", "updated_at": "2024-09-13T16:40:04.561062", "version_id": 1, "author_type": "User", "description": "string"}
980b44f4-8261-48b2-b795-512a9d8ad264	ce4a07a1-1f2e-4bc6-a635-92090d17345f	2	{"id": "ce4a07a1-1f2e-4bc6-a635-92090d17345f", "name": "tender2 bid 5", "author_id": "a6cac94a-4023-4719-b7fe-5ef7205cacd7", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:40:04.561062", "updated_at": "2024-09-13T16:40:04.561062", "version_id": 1, "author_type": "User", "description": "string"}
2538d186-375d-4505-9f2b-d6a633cf9a10	ce4a07a1-1f2e-4bc6-a635-92090d17345f	3	{"id": "ce4a07a1-1f2e-4bc6-a635-92090d17345f", "name": "version 2", "author_id": "a6cac94a-4023-4719-b7fe-5ef7205cacd7", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Created", "created_at": "2024-09-13T16:40:04.561062", "updated_at": "2024-09-13T16:40:22.132803", "version_id": 2, "author_type": "User", "description": "some description"}
080a2e08-4b26-428b-b1cd-bfff3d63ca17	ce4a07a1-1f2e-4bc6-a635-92090d17345f	4	{"id": "ce4a07a1-1f2e-4bc6-a635-92090d17345f", "name": "version 2", "author_id": "a6cac94a-4023-4719-b7fe-5ef7205cacd7", "tender_id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "bid_status": "Published", "created_at": "2024-09-13T16:40:04.561062", "updated_at": "2024-09-13T16:40:33.815119", "version_id": 3, "author_type": "User", "description": "some description"}
\.


--
-- Data for Name: employee; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.employee (id, username, first_name, last_name, created_at, updated_at) FROM stdin;
a6cac94a-4023-4719-b7fe-5ef7205cacd7	test_user	bob	builder	2024-09-10 16:46:05.475043	2024-09-10 16:46:05.475043
\.


--
-- Data for Name: organization; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.organization (id, name, description, type, created_at, updated_at) FROM stdin;
37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	avito	\N	IE	2024-09-10 16:52:44.133566	2024-09-10 16:52:44.133566
\.


--
-- Data for Name: organization_responsible; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.organization_responsible (id, organization_id, user_id) FROM stdin;
9c16c309-0d7a-491b-a286-3083cafbae45	37e1fd6c-9863-4052-8dd0-c4cc00bcb25d	a6cac94a-4023-4719-b7fe-5ef7205cacd7
\.


--
-- Data for Name: tender; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.tender (id, name, description, service_type, version_id, organization_responsible_id, created_at, updated_at, moderation_status) FROM stdin;
e2bd597a-144c-4bac-8be7-5cb34b47690b	ooooo	Нужно доставить оборудовоние	Delivery	1	9c16c309-0d7a-491b-a286-3083cafbae45	2024-09-13 16:35:19.04529	2024-09-13 16:35:19.04529	Created
97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	2	Нужно доставить оборудовоние	Delivery	1	9c16c309-0d7a-491b-a286-3083cafbae45	2024-09-13 16:35:24.708722	2024-09-13 16:35:24.708722	Created
911806b3-7655-4cde-9ce6-80083df42281	3	Нужно доставить оборудовоние	Delivery	1	9c16c309-0d7a-491b-a286-3083cafbae45	2024-09-13 16:35:27.519701	2024-09-13 16:35:27.519701	Created
8ad46627-f6bb-4bc8-968a-14c6a7a8b99b	ooooo	Нужно доставить оборудовоние	Construction	2	9c16c309-0d7a-491b-a286-3083cafbae45	2024-09-13 16:35:11.705696	2024-09-13 16:35:51.673889	Created
\.


--
-- Data for Name: tender_version; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.tender_version (id, tender_id, version_id, value) FROM stdin;
9b760b86-fe3d-450e-ae74-ed518bc4139b	8ad46627-f6bb-4bc8-968a-14c6a7a8b99b	1	{"id": "8ad46627-f6bb-4bc8-968a-14c6a7a8b99b", "name": "ooooo", "created_at": "2024-09-13T16:35:11.705696", "updated_at": "2024-09-13T16:35:11.705696", "version_id": 1, "description": "Нужно доставить оборудовоние", "service_type": "Delivery", "moderation_status": "Created", "organization_responsible_id": "9c16c309-0d7a-491b-a286-3083cafbae45"}
0e624f4f-fd18-406a-9108-25089e09dc10	e2bd597a-144c-4bac-8be7-5cb34b47690b	1	{"id": "e2bd597a-144c-4bac-8be7-5cb34b47690b", "name": "ooooo", "created_at": "2024-09-13T16:35:19.04529", "updated_at": "2024-09-13T16:35:19.04529", "version_id": 1, "description": "Нужно доставить оборудовоние", "service_type": "Delivery", "moderation_status": "Created", "organization_responsible_id": "9c16c309-0d7a-491b-a286-3083cafbae45"}
4b9f38da-90cd-4d14-a1aa-801f66d17941	97f0adb7-39cb-4048-adbd-43b0d2d7ea4f	1	{"id": "97f0adb7-39cb-4048-adbd-43b0d2d7ea4f", "name": "2", "created_at": "2024-09-13T16:35:24.708722", "updated_at": "2024-09-13T16:35:24.708722", "version_id": 1, "description": "Нужно доставить оборудовоние", "service_type": "Delivery", "moderation_status": "Created", "organization_responsible_id": "9c16c309-0d7a-491b-a286-3083cafbae45"}
9d18edc7-ef23-46cb-b35a-8d77f40b1ebf	911806b3-7655-4cde-9ce6-80083df42281	1	{"id": "911806b3-7655-4cde-9ce6-80083df42281", "name": "3", "created_at": "2024-09-13T16:35:27.519701", "updated_at": "2024-09-13T16:35:27.519701", "version_id": 1, "description": "Нужно доставить оборудовоние", "service_type": "Delivery", "moderation_status": "Created", "organization_responsible_id": "9c16c309-0d7a-491b-a286-3083cafbae45"}
49df316b-134d-40cd-84aa-75a3358c87b6	8ad46627-f6bb-4bc8-968a-14c6a7a8b99b	2	{"id": "8ad46627-f6bb-4bc8-968a-14c6a7a8b99b", "name": "ooooo", "created_at": "2024-09-13T16:35:11.705696", "updated_at": "2024-09-13T16:35:11.705696", "version_id": 1, "description": "Нужно доставить оборудовоние", "service_type": "Delivery", "moderation_status": "Created", "organization_responsible_id": "9c16c309-0d7a-491b-a286-3083cafbae45"}
\.


--
-- Name: bid_decisions bid_desitions_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_decisions
    ADD CONSTRAINT bid_desitions_pkey PRIMARY KEY (id);


--
-- Name: bid_feedback bid_feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_feedback
    ADD CONSTRAINT bid_feedback_pkey PRIMARY KEY (id);


--
-- Name: bid bid_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid
    ADD CONSTRAINT bid_pkey PRIMARY KEY (id);


--
-- Name: bid_version bid_version_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_version
    ADD CONSTRAINT bid_version_pkey PRIMARY KEY (id);


--
-- Name: employee employee_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_pkey PRIMARY KEY (id);


--
-- Name: employee employee_username_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_username_key UNIQUE (username);


--
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- Name: organization_responsible organization_responsible_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_pkey PRIMARY KEY (id);


--
-- Name: tender tender_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.tender
    ADD CONSTRAINT tender_pkey PRIMARY KEY (id);


--
-- Name: tender_version tender_version_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.tender_version
    ADD CONSTRAINT tender_version_pkey PRIMARY KEY (id);


--
-- Name: bid bids_change_trigger; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER bids_change_trigger AFTER INSERT OR UPDATE ON public.bid FOR EACH ROW EXECUTE FUNCTION public.log_bid_changes();


--
-- Name: tender tender_change_trigger; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER tender_change_trigger AFTER INSERT OR UPDATE ON public.tender FOR EACH ROW EXECUTE FUNCTION public.log_tender_changes();


--
-- Name: bid_decisions bid_desitions_bid_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_decisions
    ADD CONSTRAINT bid_desitions_bid_id_fkey FOREIGN KEY (bid_id) REFERENCES public.bid(id) ON DELETE CASCADE;


--
-- Name: bid_feedback bid_feedback_bid_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_feedback
    ADD CONSTRAINT bid_feedback_bid_id_fkey FOREIGN KEY (bid_id) REFERENCES public.bid(id) ON DELETE CASCADE;


--
-- Name: bid bid_tender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid
    ADD CONSTRAINT bid_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES public.tender(id);


--
-- Name: bid_version bid_version_bid_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bid_version
    ADD CONSTRAINT bid_version_bid_id_fkey FOREIGN KEY (bid_id) REFERENCES public.bid(id) ON DELETE CASCADE;


--
-- Name: organization_responsible organization_responsible_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- Name: organization_responsible organization_responsible_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.organization_responsible
    ADD CONSTRAINT organization_responsible_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.employee(id) ON DELETE CASCADE;


--
-- Name: tender tender_organization_responsible_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.tender
    ADD CONSTRAINT tender_organization_responsible_id_fkey FOREIGN KEY (organization_responsible_id) REFERENCES public.organization_responsible(id);


--
-- Name: tender_version tender_version_tender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.tender_version
    ADD CONSTRAINT tender_version_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES public.tender(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

