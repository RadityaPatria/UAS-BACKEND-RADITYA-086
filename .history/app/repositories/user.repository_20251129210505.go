package repositories

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)


// =====================================================
// GET USER BY USERNAME ATAU EMAIL (untuk Login)
// =====================================================
func GetUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, username, email, password_hash, full_name, role_id, 
		        is_active, created_at, updated_at
		 FROM users 
		 WHERE username=$1 OR email=$1`,
		identifier,
	)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.roleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}


// =====================================================
// GET USER BY ID
// =====================================================
func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, username, email, password_hash, full_name, role_id,
		        is_active, created_at, updated_at
		 FROM users WHERE id=$1`,
		id,
	)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.roleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}


// =====================================================
// CREATE USER (ADMIN)
// =====================================================
func CreateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash, full_name,
		                    role_id, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.FullName, u.RoleID, u.IsActive,
	)
	return err
}


// =====================================================
// UPDATE USER (ADMIN)
// =====================================================
func UpdateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE users 
		 SET username=$1, email=$2, full_name=$3, role_id=$4,
		     is_active=$5, updated_at=NOW()
		 WHERE id=$6`,
		u.Username, u.Email, u.FullName, u.RoleID,
		u.IsActive, u.ID,
	)

	return err
}


// =====================================================
// DELETE USER (ADMIN)
// =====================================================
func DeleteUser(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM users WHERE id=$1`, id)
	return err
}


// =====================================================
// GET USER PERMISSIONS (untuk JWT & RBAC)
// =====================================================
func GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT p.name
		 FROM role_permissions rp
		 JOIN permissions p ON rp.permission_id = p.id
		 JOIN users u ON u.role_id = rp.role_id
		 WHERE u.id = $1
		 ORDER BY p.name`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string

	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		list = append(list, perm)
	}

	return list, nil
}


// =====================================================
// GET ROLE NAME BY ID
// =====================================================
func GetRoleNameByID(ctx context.Context, roleID string) (string, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT name FROM roles WHERE id=$1`, roleID)

	var name string
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
