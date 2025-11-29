package repositories

import (
    "context"
    "UAS-backend/database"
    "UAS-backend/app/models"
    "github.com/google/uuid"
)

// =====================================
// AssignPermissionToRole()
// Admin memberikan permission ke role
// =====================================
func AssignPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
    query := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
    `
    _, err := database.DB.Exec(ctx, query, roleID, permissionID)
    return err
}

// =====================================
// GetPermissionsByRole()
// Ambil semua permission milik role
// dipakai saat login (buat JWT)
// =====================================
func GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
    query := `
        SELECT p.id, p.name, p.resource, p.action, p.description
        FROM permissions p
        JOIN role_permissions rp ON rp.permission_id = p.id
        WHERE rp.role_id = $1
    `
    rows, err := database.DB.Query(ctx, query, roleID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var perms []models.Permission

    for rows.Next() {
        var p models.Permission
        err := rows.Scan(
            &p.ID, &p.Name, &p.Resource, &p.Action, &p.Description,
        )
        if err != nil {
            return nil, err
        }
        perms = append(perms, p)
    }

    return perms, nil
}

// =====================================
// RemovePermissionFromRole()
// Admin mencabut permission
// =====================================
func RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
    query := `
        DELETE FROM role_permissions 
        WHERE role_id = $1 AND permission_id = $2
    `
    _, err := database.DB.Exec(ctx, query, roleID, permissionID)
    return err
}
