package repositories

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/database"
)

// GetAllUsers -> ambil semua user | FR-009
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash,
			&u.FullName, &u.RoleID, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByIdentifier -> login pakai username / email | FR-001
func GetUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, token_version, created_at, updated_at
		FROM users
		WHERE username=$1 OR email=$1`, identifier)

	var u models.User
	if err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive,
		&u.TokenVersion, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByID -> validasi JWT & token_version | FR-002
func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, token_version, created_at, updated_at
		FROM users WHERE id=$1`, id)

	var u models.User
	if err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive,
		&u.TokenVersion, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

// IncrementTokenVersion -> logout (invalidate token) | FR-001
func IncrementTokenVersion(ctx context.Context, userID string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET token_version = token_version + 1 WHERE id=$1`,
		userID,
	)
	return err
}

// UpdateTokenVersion -> refresh token | FR-001
func UpdateTokenVersion(ctx context.Context, userID string, version int) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET token_version=$1, updated_at=NOW() WHERE id=$2`,
		version, userID,
	)
	return err
}

// CreateUser -> tambah user baru | FR-009
func CreateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name,
			role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.FullName, u.RoleID, u.IsActive,
	)
	return err
}

// UpdateUser -> update data user | FR-009
func UpdateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE users
		SET username=$1, email=$2, full_name=$3,
		    role_id=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6`,
		u.Username, u.Email, u.FullName,
		u.RoleID, u.IsActive, u.ID,
	)
	return err
}

// DeleteUser -> hapus user | FR-009
func DeleteUser(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx,
		`DELETE FROM users WHERE id=$1`, id)
	return err
}
