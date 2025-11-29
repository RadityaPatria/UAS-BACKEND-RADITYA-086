package repositories

import (
	"context"
	"errors"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

// Struktur repository
type UserRepository struct{}

// Instansiasi global
var UserRepo = UserRepository{}

// Ambil user by username atau email (untuk login)
func (r UserRepository) GetByUsernameOrEmail(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $1
	`
	row := database.DB.QueryRow(ctx, query, username)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email,
		&u.PasswordHash, &u.FullName,
		&u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Ambil user berdasarkan ID
func (r UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	row := database.DB.QueryRow(ctx, query, id)

	var u models.User
	if err := row.Scan(
		&u.ID, &u.Username, &u.Email,
		&u.PasswordHash, &u.FullName,
		&u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

// Buat user baru (admin)
func (r UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := database.DB.Exec(ctx, query,
		user.ID, user.Username, user.Email,
		user.PasswordHash, user.FullName,
		user.RoleID, user.IsActive,
	)
	return err
}

// Update user
func (r UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username=$2, email=$3, full_name=$4, role_id=$5, is_active=$6
		WHERE id=$1
	`
	_, err := database.DB.Exec(ctx, query,
		user.ID, user.Username, user.Email,
		user.FullName, user.RoleID, user.IsActive,
	)
	return err
}

// Delete user
func (r UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return err
}
