package billing

import (
	"avito-test/internal/models"
)

func (r *repository) GetUserByUsername(userName string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT id, first_name, last_name, created_at, updated_at
		FROM employee
			WHERE username = $1
	`
	row := r.db.QueryRow(query, &userName)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *repository) GetUserById(uuid string) (models.User, error) {
	user := models.User{}
	query := `
		SELECT id, first_name, last_name, created_at, updated_at
		FROM employee
			WHERE id = $1
	`
	row := r.db.QueryRow(query, &uuid)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return models.User{}, err
	}
	return user, nil
}
