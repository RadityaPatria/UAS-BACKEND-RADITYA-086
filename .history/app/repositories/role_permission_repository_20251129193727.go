package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// GetPermissionIDsByRoleID → ambil permission_id dari role_id
func GetPermissionIDsByRoleID(ctx context.Context, roleID string) ([]string, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT permission_id FROM role_permissions WHERE role_id = $1`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		ids = append(ids, pid)
	}
	return ids, nil
}

// AssignPermissionToRole → assign permission ke role
func AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO role_permissions (role_id, permission_id)
		 VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		roleID, permissionID)
	return err
}

// RemovePermissionFromRole → hapus permission dari role
func RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`,
		roleID, permissionID)
	return err
}

// GetAllRolePermissions → ambil semua role-permission
func GetAllRolePermissions(ctx context.Context) ([]models.RolePermission, error) {
	rows, err := database.DB.Query(ctx, `SELECT role_id, permission_id FROM role_permissions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.RolePermission
	for rows.Next() {
		var rp models.RolePermission
		if err := rows.Scan(&rp.RoleID, &rp.PermissionID); err != nil {
			return nil, err
		}
		list = append(list, rp)
	}
	return list, nil
}
