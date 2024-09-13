package billing

import (
	"avito-test/internal/models"
	"avito-test/internal/storage/billing/convertion"
	"database/sql"
	"fmt"
)

func (r *repository) CreateTender(name, description, serviceType, moderationStatus, responsibleId string) (string, string, error) {
	var id, date string
	query :=
		`
		INSERT INTO 
			tender(name, description, service_type, moderation_status, organization_responsible_id)
			VALUES($1,$2,$3,$4,$5)
			RETURNING id, created_at
	
	`
	row := r.db.QueryRow(query, &name, &description, &serviceType, &moderationStatus, &responsibleId)
	if err := row.Scan(&id, &date); err != nil {
		return "", "", err
	}
	return id, date, nil
}
func (r *repository) UpdateTender(tenderId string, builder models.TenderBuilder) (models.Tender, error) {
	var (
		tender          = models.Tender{}
		service, status string
		err             error
	)
	query :=
		`
		UPDATE tender 
			SET
				name = COALESCE(NULLIF($1, ''), name),
				description = COALESCE(NULLIF($2, ''), description),
				service_type = COALESCE(NULLIF($3, '')::service_type, service_type),
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			WHERE id = $4
			RETURNING id, name, description, moderation_status, service_type, version_id, created_at
	`
	if builder.ServiceType != nil {
		service = builder.ServiceType.String()
	}
	if builder.Status != nil {
		status = builder.Status.String()
	}
	row := r.db.QueryRow(query, &builder.Name, &builder.Description, &service, &tenderId)
	if err = row.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service,
		&tender.Version, &tender.CreatedAt); err != nil {
		return models.Tender{}, err
	}
	tender.Status, err = convertion.TenderStatusFromString(status)
	if err != nil {
		return models.Tender{}, err
	}
	tender.ServiceType, err = convertion.ServiceFromString(service)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (r *repository) SetTenderStatus(tenderId string, value models.TenderStatus) (models.Tender, error) {
	var (
		tender          = models.Tender{}
		service, status string
		err             error
	)
	query :=
		`
		UPDATE tender 
			SET
				moderation_status = $1,
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			WHERE id = $2
			RETURNING id, name, description, moderation_status, service_type, version_id, created_at
	`
	status = value.String()
	row := r.db.QueryRow(query, &status, &tenderId)
	if err = row.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service,
		&tender.Version, &tender.CreatedAt); err != nil {
		return models.Tender{}, err
	}

	tender.Status = value
	tender.ServiceType, err = convertion.ServiceFromString(service)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (r *repository) GetTender(tenderId string) (models.Tender, error) {
	var (
		tender          models.Tender
		service, status string
		err             error
	)
	query :=
		`
		SELECT id, name, description, moderation_status, service_type, version_id, created_at
		FROM tender
		WHERE id = $1
	`
	row := r.db.QueryRow(query, &tenderId)
	if err := row.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service, &tender.Version, &tender.CreatedAt); err != nil {
		return models.Tender{}, err
	}
	tender.ServiceType, err = convertion.ServiceFromString(service)
	if err != nil {
		return models.Tender{}, err
	}
	tender.Status, err = convertion.TenderStatusFromString(status)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (r *repository) RollbackTender(tenderId string, version int) (models.Tender, error) {
	var (
		tender          = models.Tender{}
		status, service string
		err             error
	)
	query :=
		`
		UPDATE tender 
			SET
				name = value->>'name',
				description = value->>'description',
				service_type = (value->>'service_type')::service_type,
				moderation_status = (value->>'moderation_status')::moderation_status,
				updated_at = CURRENT_TIMESTAMP,
				version_id = version_id + 1 
			FROM (
				SELECT (value)::jsonb
				FROM tender_version
				WHERE tender_id = $1 AND version_id = $2
			) as value
			WHERE id = $1
			RETURNING id, name, description, moderation_status, service_type, version_id, created_at
	`

	row := r.db.QueryRow(query, &tenderId, &version)
	if err = row.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service,
		&tender.Version, &tender.CreatedAt); err != nil {
		return models.Tender{}, err
	}

	tender.Status, err = convertion.TenderStatusFromString(status)
	if err != nil {
		return models.Tender{}, err
	}
	tender.ServiceType, err = convertion.ServiceFromString(service)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (r *repository) GetTenders(limit, offset int, serviceTypes []models.EnumService) ([]models.Tender, error) {
	var (
		result          = []models.Tender{}
		tender          = models.Tender{}
		status, service string
		err             error
	)
	services := ""
	for _, service := range serviceTypes {
		services += fmt.Sprintf(`'%s',`, service.String())
	}
	services = services[:len(services)-1]
	query := fmt.Sprintf(
		`
		SELECT id, name, description, moderation_status, service_type, version_id, created_at
		FROM tender
		WHERE service_type in(%s)
		ORDER BY name DESC
		LIMIT $1
		OFFSET $2
	`, services)
	rows, err := r.db.Query(query, &limit, &offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service,
			&tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, err
		}
		tender.Status, err = convertion.TenderStatusFromString(status)
		if err != nil {
			return nil, err
		}
		tender.ServiceType, err = convertion.ServiceFromString(service)
		if err != nil {
			return nil, err
		}
		result = append(result, tender)
	}
	return result, nil
}
func (r *repository) GetTenderForBid(bidId string) (models.Tender, error) {
	var (
		tender          models.Tender
		service, status string
		err             error
	)
	query :=
		`
		SELECT tender.id, tender.name, tender.description, tender.moderation_status, tender.service_type, tender.version_id, tender.created_at
		FROM tender INNER JOIN bid on (tender.id = bid.tender_id)
		WHERE bid.id = $1
	`
	row := r.db.QueryRow(query, &bidId)
	if err := row.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service, &tender.Version, &tender.CreatedAt); err != nil {
		return models.Tender{}, err
	}
	tender.ServiceType, err = convertion.ServiceFromString(service)
	if err != nil {
		return models.Tender{}, err
	}
	tender.Status, err = convertion.TenderStatusFromString(status)
	if err != nil {
		return models.Tender{}, err
	}
	return tender, nil
}

func (r *repository) GetUserTenders(limit, offset int, userId string) ([]models.Tender, error) {
	var (
		result          = []models.Tender{}
		tender          models.Tender
		service, status string
		err             error
	)
	query :=
		`
		SELECT tender.id, name, description, moderation_status, service_type, version_id, created_at
		FROM tender INNER JOIN organization_responsible ON (organization_responsible.id = tender.organization_responsible_id)
		WHERE organization_responsible.user_id = $1
		LIMIT $2 OFFSET $3
	
	`
	rows, err := r.db.Query(query, &userId, &limit, &offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &status, &service,
			&tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, err
		}
		tender.Status, err = convertion.TenderStatusFromString(status)
		if err != nil {
			return nil, err
		}
		tender.ServiceType, err = convertion.ServiceFromString(service)
		if err != nil {
			return nil, err
		}
		result = append(result, tender)
	}
	return result, nil
}
func (r *repository) ValidTenderRights(userId, tenderId string) (string, error) {
	var id string
	query := `
		SELECT organization_responsible.id
		FROM employee 
			INNER JOIN organization_responsible  ON (employee.id = organization_responsible.user_id )
			INNER JOIN tender ON (organization_responsible.id = tender.organization_responsible_id)
		WHERE employee.id = $1
			AND tender.id = $2
		LIMIT 1
	`
	row := r.db.QueryRow(query, &userId, &tenderId)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return id, nil
}
