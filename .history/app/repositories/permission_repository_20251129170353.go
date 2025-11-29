package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type PermissionRepository struct{}

var PermissionRepo = &PermissionRepository{}

// Mendapatkan semua permissions
func (r *PermissionRepository) GetAll(ctx context.Context) ([]models.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
	`
	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission

	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// Ambil permission berdasarkan ID
func (r *PermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE id=$1
	`
	row := database.DB.QueryRow(ctx, query, id)

	var p models.Permission
	if err := row.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
		return nil, err
	}

	return &p, nil
}

// Buat permission baru (admin)
func (r *PermissionRepository) Create(ctx context.Context, p *models.Permission) error {
	p.ID = uuid.New()

	_, err := database.DB.Exec(ctx,
		`INSERT INTO permissions (id, name, resource, action, description)
		 VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.Name, p.Resource, p.Action, p.Description,
	)

	return err
}

// Update permission (admin)
func (r *PermissionRepository) Update(ctx context.Context, p *models.Permission) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE permissions
		 SET name=$1, resource=$2, action=$3, description=$4
		 WHERE id=$5`,
		p.Name, p.Resource, p.Action, p.Description,
		p.ID,
	)

	return err
}

// Hapus permission (admin)
func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM permissions WHERE id=$1`, id)
	return err
}

// Ambil permission berdasarkan Role
func (r *PermissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	query := `
		SELECT p.name 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := database.DB.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}

	return perms, nil
}
