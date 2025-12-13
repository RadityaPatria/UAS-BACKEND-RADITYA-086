package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// GetAllRoles -> ambil semua role | FR-002
func GetAllRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, description, created_at
		 FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

// GetRoleByID -> ambil role berdasarkan id | FR-002
func GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, description, created_at
		 FROM roles
		 WHERE id=$1`, id)

	var r models.Role
	if err := row.Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &r, nil
}

// GetRoleByName -> ambil role berdasarkan nama | FR-002
func GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, description, created_at
		 FROM roles
		 WHERE name=$1`, name)

	var r models.Role
	if err := row.Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &r, nil
}

// CreateRole -> tambah role baru | FR-002
func CreateRole(ctx context.Context, role *models.Role) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO roles (id, name, description, created_at)
		 VALUES ($1,$2,$3,NOW())`,
		role.ID, role.Name, role.Description,
	)
	return err
}

// UpdateRole -> ubah data role | FR-002
func UpdateRole(ctx context.Context, role *models.Role) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE roles
		 SET name=$1, description=$2
		 WHERE id=$3`,
		role.Name, role.Description, role.ID,
	)
	return err
}

// DeleteRole -> hapus role | FR-002
func DeleteRole(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM roles WHERE id=$1`, id)
	return err
}
