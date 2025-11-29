package repositories

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"github.com/google/uuid"
)

type UserRepository struct{}

var UserRepo = &UserRepository{}

// Fungsi ini untuk mencari user saat login (username OR email)
func (r *UserRepository) FindByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $1
	`

	row := database.DB.QueryRow(ctx, query, identifier)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Mendapatkan user berdasarkan ID (untuk profile)
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Membuat user (dipakai oleh admin)
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := database.DB.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// Update user oleh admin
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	_, err := database.DB.Exec(ctx, `
		UPDATE users SET 
			username=$1,
			email=$2,
			full_name=$3,
			role_id=$4,
			is_active=$5,
			updated_at=$6
		WHERE id=$7
	`,
		user.Username,
		user.Email,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// Hard delete user (admin)
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
