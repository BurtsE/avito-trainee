package billing

import "avito-test/internal/models"

func (r *repository) GetOrganizationById(uuid string) (models.DummyOrganization, error) {
	organization := models.DummyOrganization{}
	query := `
		SELECT id, name, COALESCE(description,'') , type, created_at, updated_at
		FROM organization 
			WHERE id = $1
	`
	row := r.db.QueryRow(query, &uuid)
	if err := row.Scan(&organization.Id, &organization.Name, &organization.Desription, &organization.Type,
		&organization.UpdatedAt, &organization.UpdatedAt); err != nil {
		return models.DummyOrganization{}, err
	}
	return organization, nil
}

func (r *repository) GetResponsibleOrganization(userId, organizationId string) (string, error) {
	var uuid string
	query := `
		SELECT id
		FROM organization_responsible 
			WHERE user_id = $1 AND organization_id =$2
	`
	row := r.db.QueryRow(query, &userId, &organizationId)
	if err := row.Scan(&uuid); err != nil {
		return "", err
	}
	return uuid, nil
}

func (r *repository) GetOrganizationId(userId string) (string, error) {
	var uuid string
	query := `
		SELECT organization_id
		FROM organization_responsible 
			WHERE user_id = $1
	`
	row := r.db.QueryRow(query, &userId, &userId)
	if err := row.Scan(&uuid); err != nil {
		return "", err
	}
	return uuid, nil
}
