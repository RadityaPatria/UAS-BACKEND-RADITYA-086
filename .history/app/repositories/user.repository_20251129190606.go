
package repository

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

// ambil user by email
func GetUserByEmail(db *sql.DB, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1
	`

	user := new(model.User)
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ambil user by id
func GetUserByID(db *sql.DB, id string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`

	user := new(model.User)
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ambil user by username
func GetUserByUsername(db *sql.DB, username string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.username = $1
	`

	user := new(model.User)
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// ambil permissions user
func GetUserPermissions(db *sql.DB, userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1
		ORDER BY p.name
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// ambil role name by id
func GetRoleNameByID(db *sql.DB, roleID string) (string, error) {
	query := `SELECT name FROM roles WHERE id = $1`

	var roleName string
	err := db.QueryRow(query, roleID).Scan(&roleName)
	if err != nil {
		return "", err
	}

	return roleName, nil
}

saya hanya nanya kenapa format reponya beda dengan ini