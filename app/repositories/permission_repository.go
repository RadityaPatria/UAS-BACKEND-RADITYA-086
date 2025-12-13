package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// GetAllPermissions -> ambil semua permission | FR-003
func GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, name, resource, action, description
		 FROM permissions
		 ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Resource,
			&p.Action,
			&p.Description,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

// GetPermissionByID -> ambil permission berdasarkan id | FR-003
func GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, name, resource, action, description
		 FROM permissions
		 WHERE id = $1`, id)

	var p models.Permission
	if err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Resource,
		&p.Action,
		&p.Description,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPermissionsByIDs -> ambil banyak permission saat login | FR-001
func GetPermissionsByIDs(ctx context.Context, ids []string) ([]models.Permission, error) {
	if len(ids) == 0 {
		return []models.Permission{}, nil
	}

	rows, err := database.DB.Query(ctx,
		`SELECT id, name, resource, action, description
		 FROM permissions
		 WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Resource,
			&p.Action,
			&p.Description,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

// CreatePermission -> tambah permission baru | FR-003
func CreatePermission(ctx context.Context, p *models.Permission) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO permissions (id, name, resource, action, description)
		 VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.Name, p.Resource, p.Action, p.Description,
	)
	return err
}

// UpdatePermission -> ubah data permission | FR-003
func UpdatePermission(ctx context.Context, p *models.Permission) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE permissions
		 SET name=$2, resource=$3, action=$4, description=$5
		 WHERE id=$1`,
		p.ID, p.Name, p.Resource, p.Action, p.Description,
	)
	return err
}

// DeletePermission -> hapus permission | FR-003
func DeletePermission(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM permissions WHERE id=$1`, id)
	return err
}
