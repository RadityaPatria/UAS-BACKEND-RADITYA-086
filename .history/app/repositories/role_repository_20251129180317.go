package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type RoleRepository struct{}

var RoleRepo = &RoleRepository{}

// Ambil semua role
func (r *RoleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
	`
	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// Ambil 1 role berdasarkan ID
func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var role models.Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
		return nil, err
	}

	return &role, nil
}

// Ambil role berdasarkan nama (dipakai admin)
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE name = $1
	`

	row := database.DB.QueryRow(ctx, query, name)

	var role models.Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
		return nil, err
	}

	return &role, nil
}

// Buat role baru (ADMIN ONLY)
func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	role.ID = uuid.New()

	_, err := database.DB.Exec(ctx,
		`INSERT INTO roles (id, name, description)
		 VALUES ($1, $2, $3)`,
		role.ID, role.Name, role.Description,
	)

	return err
}

// Update role
func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE roles SET name=$1, description=$2 WHERE id=$3`,
		role.Name, role.Description, role.ID,
	)

	return err
}

// Hapus role
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM roles WHERE id=$1`,
		id,
	)

	return err
}
