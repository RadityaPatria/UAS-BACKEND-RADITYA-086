package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type RolePermissionRepository struct{}

var RolePermissionRepo = &RolePermissionRepository{}

// Ambil semua permission_id berdasarkan role_id
func (r *RolePermissionRepository) GetPermissionIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT permission_id
		FROM role_permissions
		WHERE role_id = $1
	`

	rows, err := database.DB.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID

	for rows.Next() {
		var pid uuid.UUID
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		ids = append(ids, pid)
	}

	return ids, nil
}

// Assign permission ke role
func (r *RolePermissionRepository) AssignPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO role_permissions (role_id, permission_id)
		 VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`,
		roleID, permissionID,
	)

	return err
}

// Hapus permission dari role
func (r *RolePermissionRepository) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM role_permissions
		 WHERE role_id=$1 AND permission_id=$2`,
		roleID, permissionID,
	)

	return err
}

// Ambil role-permission lengkap (admin)
func (r *RolePermissionRepository) GetAll(ctx context.Context) ([]models.RolePermission, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT role_id, permission_id FROM role_permissions`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.RolePermission

	for rows.Next() {
		var rp models.RolePermission
		if err := rows.Scan(&rp.RoleID, &rp.PermissionID); err != nil {
			return nil, err
		}
		data = append(data, rp)
	}

	return data, nil
}
