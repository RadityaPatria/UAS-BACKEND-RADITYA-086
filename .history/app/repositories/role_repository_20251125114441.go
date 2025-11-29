package repositories

import (
    "context"
    "UAS-backend/database"
    "UAS-backend/app/models"
    "github.com/google/uuid"
)

// =====================================
// CreateRole()
// Dipakai admin untuk membuat role baru
// =====================================
func CreateRole(ctx context.Context, r *models.Role) error {
    query := `
        INSERT INTO roles (id, name, description, created_at)
        VALUES ($1, $2, $3, NOW())
    `
    _, err := database.DB.Exec(ctx, query, r.ID, r.Name, r.Description)
    return err
}

// =====================================
// GetRoleByID()
// Ambil 1 role berdasarkan ID
// =====================================
func GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
    query := `
        SELECT id, name, description, created_at
        FROM roles
        WHERE id = $1
    `
    row := database.DB.QueryRow(ctx, query, id)

    var role models.Role
    err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
    if err != nil {
        return nil, err
    }

    return &role, nil
}

// =====================================
// GetAllRoles()
// Admin lihat semua role
// =====================================
func GetAllRoles(ctx context.Context) ([]models.Role, error) {
    query := `
        SELECT id, name, description, created_at
        FROM roles ORDER BY created_at DESC
    `
    rows, err := database.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var roles []models.Role

    for rows.Next() {
        var r models.Role
        err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt)
        if err != nil {
            return nil, err
        }

        roles = append(roles, r)
    }

    return roles, nil
}

// =====================================
// UpdateRole()
// Admin edit role
// =====================================
func UpdateRole(ctx context.Context, r *models.Role) error {
    query := `
        UPDATE roles SET name = $2, description = $3
        WHERE id = $1
    `
    _, err := database.DB.Exec(ctx, query, r.ID, r.Name, r.Description)
    return err
}

// =====================================
// DeleteRole()
// Admin hapus role
// =====================================
func DeleteRole(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM roles WHERE id = $1`
    _, err := database.DB.Exec(ctx, query, id)
    return err
}
