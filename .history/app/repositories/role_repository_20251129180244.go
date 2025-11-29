package repositories

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

// ===============================
// GET all roles
// ===============================
func GetAllRoles(ctx context.Context) ([]models.Role, error) {

	rows, err := database.DB.Query(ctx,
		`SELECT id, name, description, created_at, updated_at FROM roles ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Role

	for rows.Next() {
		var r models.Role
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}


// ===============================
// GET role by ID
// ===============================
func GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {

	row := database.DB.QueryRow(ctx,
		`SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`,
		id,
	)

	var r models.Role
	err := row.Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
		&r.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}


// ===============================
// CREATE role
// ===============================
func CreateRole(ctx context.Context, role *models.Role) error {

	role.ID = uuid.New()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := database.DB.Exec(ctx,
		`INSERT INTO roles (id, name, description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	return err
}


// ===============================
// UPDATE role
// ===============================
func UpdateRole(ctx context.Context, role *models.Role) error {

	role.UpdatedAt = time.Now()

	_, err := database.DB.Exec(ctx,
		`UPDATE roles SET 
			name = $1,
			description = $2,
			updated_at = $3
		WHERE id = $4`,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)

	return err
}


// ===============================
// DELETE role
// ===============================
func DeleteRole(ctx context.Context, id uuid.UUID) error {

	res, err := database.DB.Exec(ctx,
		`DELETE FROM roles WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("role not found")
	}

	return nil
}
