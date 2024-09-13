package billing

import (
	"avito-test/internal/models"
	"avito-test/internal/storage/billing/convertion"
	"database/sql"
)

func (r *repository) GetBid(tenderId string) (models.Bid, error) {
	var (
		bid                models.Bid
		status, authorType string
		err                error
	)
	query :=
		`
		SELECT id, name, bid_status, author_type, description, author_id, version_id, created_at
		FROM bid
		WHERE id = $1
	`
	row := r.db.QueryRow(query, &tenderId)
	if err := row.Scan(&bid.Id, &bid.Name, &status, &authorType, &bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt); err != nil {
		return models.Bid{}, err
	}
	bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
	if err != nil {
		return models.Bid{}, err
	}
	bid.Status, err = convertion.BidStatusFromString(status)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (r *repository) CreateBid(builder models.BidsBuilder) (models.Bid, error) {
	var (
		bid                models.Bid
		authorType, status string
		err                error
	)
	query :=
		`
		INSERT INTO 
			bid(name, description, tender_id, author_type, author_id)
			VALUES($1,$2,$3,$4,$5)
			RETURNING id, name, bid_status, author_type, description, author_id, version_id, created_at
	
	`
	authorType = builder.AuthorType.String()
	row := r.db.QueryRow(query, &builder.Name, &builder.Description, &builder.TenderId, &authorType, &builder.AuthorId)
	if err := row.Scan(&bid.Id, &bid.Name, &status, &authorType,
		&bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt); err != nil {
		return models.Bid{}, err
	}
	bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
	if err != nil {
		return models.Bid{}, err
	}
	bid.Status, err = convertion.BidStatusFromString(status)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (r *repository) GetUserBids(limit, offset int, userId string) ([]models.Bid, error) {
	var (
		result             = []models.Bid{}
		bid                models.Bid
		status, authorType string
		err                error
	)
	query :=
		`
		SELECT bid.id, name, bid_status, author_type, description, author_id, version_id, created_at
		FROM bid 
		WHERE author_id =$1
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, &userId, &limit, &offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&bid.Id, &bid.Name, &status, &authorType, &bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
		if err != nil {
			return nil, err
		}
		bid.Status, err = convertion.BidStatusFromString(status)
		if err != nil {
			return nil, err
		}
		result = append(result, bid)
	}
	return result, nil
}

func (r *repository) UpdateBid(bidId string, builder models.BidsBuilder) (models.Bid, error) {
	var (
		bid                = models.Bid{}
		status, authorType string
		err                error
	)
	query :=
		`
		UPDATE bid 
			SET
				name = COALESCE(NULLIF($1, ''), name),
				description = COALESCE(NULLIF($2, ''), description),
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			WHERE id = $3
			RETURNING id, name, bid_status, author_type, description, author_id, version_id, created_at
	`

	row := r.db.QueryRow(query, &builder.Name, &builder.Description, &bidId)
	if err = row.Scan(&bid.Id, &bid.Name, &status, &authorType, &bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt); err != nil {
		return models.Bid{}, err
	}
	bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
	if err != nil {
		return models.Bid{}, err
	}
	bid.Status, err = convertion.BidStatusFromString(status)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (r *repository) SetBidStatus(bidId string, value models.BidStatus) (models.Bid, error) {
	var (
		bid                = models.Bid{}
		status, authorType string
		err                error
	)
	query :=
		`
		UPDATE bid 
			SET
				bid_status = $1,
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			WHERE id = $2
			RETURNING id, name, bid_status, author_type, description, author_id, version_id, created_at
	`
	status = value.String()
	row := r.db.QueryRow(query, &status, &bidId)
	if err = row.Scan(&bid.Id, &bid.Name, &status, &authorType,
		&bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt); err != nil {
		return models.Bid{}, err
	}

	bid.Status = value
	bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (r *repository) RollbackBid(bidId string, version int) (models.Bid, error) {
	var (
		bid                = models.Bid{}
		status, authorType string
		err                error
	)
	query :=
		`
		UPDATE bid 
			SET
				name = value->>'name',
				description = value->>'description',
				bid_status = (value->>'bid_status')::bid_status,
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			FROM (
				SELECT (value)::jsonb
				FROM bid_version
				WHERE bid_id = $1 AND version_id = $2
			) as value
			WHERE id = $1
			RETURNING id, name, bid_status, author_type, description, author_id, version_id, created_at
	`

	row := r.db.QueryRow(query, &bidId, &version)
	if err = row.Scan(&bid.Id, &bid.Name, &status, &authorType,
		&bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt); err != nil {
		return models.Bid{}, err
	}

	bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
	if err != nil {
		return models.Bid{}, err
	}
	bid.Status, err = convertion.BidStatusFromString(status)
	if err != nil {
		return models.Bid{}, err
	}
	return bid, nil
}

func (r *repository) GetBidsForTender(tenderId string, limit, offset int) ([]models.Bid, error) {
	var (
		result             = []models.Bid{}
		bid                = models.Bid{}
		status, authorType string
		err                error
	)
	query :=
		`
		SELECT id, name, bid_status, author_type, description, author_id, version_id, created_at
		FROM bid
		WHERE tender_id = $1
		ORDER BY name DESC
		LIMIT $2
		OFFSET $3
	`
	rows, err := r.db.Query(query, &tenderId, &limit, &offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&bid.Id, &bid.Name, &status, &authorType,
			&bid.Description, &bid.AuthorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bid.AuthorType, err = convertion.AuthorTypeFromString(authorType)
		if err != nil {
			return nil, err
		}
		bid.Status, err = convertion.BidStatusFromString(status)
		if err != nil {
			return nil, err
		}
		result = append(result, bid)
	}
	return result, nil
}

func (r *repository) ValidBidRights(authorId, bidId string) (string, error) {
	var id string
	query := `
		WITH org as(
			SELECT id, user_id, organization_id
			FROM organization_responsible
			WHERE (organization_responsible.user_id = $1 OR organization_responsible.organization_id = $1)
			)
		SELECT org.id
		FROM bid CROSS JOIN org
		WHERE bid.id = $2 and (org.user_id = $1 OR org.organization_id = $1)
		LIMIT 1
	`
	row := r.db.QueryRow(query, &authorId, &bidId)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *repository) AddBidFeedback(bidId, feedback string) error {
	query :=
		`
		INSERT INTO bid_feedback(bid_id, value)
		VALUES($1,$2)
	`
	_, err := r.db.Exec(query, &bidId, feedback)
	if err != nil {
		return err
	}
	return nil
}
