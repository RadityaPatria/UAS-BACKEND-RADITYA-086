package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type PermissionRepository struct{}

var PermissionRepo = &PermissionRepository{}

// ----------------------------------------------------
// Ambil Semua Permission
// ----------------------------------------------------
func (r *PermissionRepository) GetAll(ctx context.Context) ([]models.Permission, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, resource, action, description 
		 FROM permissions`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission

	for rows.Next() {
		var p models.Permission
		err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// ----------------------------------------------------
// Ambil Permission berdasarkan ID
// ----------------------------------------------------
func (r *PermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, resource, action, description
		 FROM permissions WHERE id=$1`,
		id,
	)

	var p models.Permission
	err := row.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ----------------------------------------------------
// Ambil permission berdasarkan daftar UUID
// digunakan saat login untuk convert UUID → string permission
// ----------------------------------------------------
func (r *PermissionRepository) GetManyByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Permission, error) {
	if len(ids) == 0 {
		return []models.Permission{}, nil
	}

	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE id = ANY($1)
	`

	rows, err := database.DB.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission

	for rows.Next() {
		var p models.Permission
		err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// ----------------------------------------------------
// Tambah Permission Baru (Admin)
// ----------------------------------------------------
func (r *PermissionRepository) Create(ctx context.Context, p *models.Permission) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO permissions (id, name, resource, action, description)
		 VALUES ($1, $2, $3, $4, $5)`,
		p.ID, p.Name, p.Resource, p.Action, p.Description,
	)

	return err
}

// ----------------------------------------------------
// Update Permission
// ----------------------------------------------------
func (r *PermissionRepository) Update(ctx context.Context, p *models.Permission) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE permissions
		 SET name=$2, resource=$3, action=$4, description=$5
		 WHERE id=$1`,
		p.ID, p.Name, p.Resource, p.Action, p.Description,
	)

	return err
}

// ----------------------------------------------------
// Hapus Permission
// ----------------------------------------------------
func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM permissions WHERE id=$1`,
		id,
	)

	return err
}
