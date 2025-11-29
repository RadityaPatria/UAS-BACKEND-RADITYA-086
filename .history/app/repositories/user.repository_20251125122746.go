package repositories

import (
	"context"
	"time"

	"UAS-backend/database"
	"UAS-backend/app/models"

	"github.com/google/uuid"
)

// UserRepository menangani semua query dan operasi database untuk tabel users.
type UserRepository struct{}

// NewUserRepository membuat instance repository.
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByUsernameOrEmail mencari user berdasarkan username ATAU email (untuk login).
func (r *UserRepository) FindByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $1
	`
	row := database.DB.QueryRow(ctx, query, identifier)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID mengambil user berdasarkan ID.
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create menambahkan user baru (admin create user, mahasiswa registrasi, dsb).
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := database.DB.Exec(
		ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.FullName, user.RoleID, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

// Update memperbarui data user (admin update user).
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=$6
		WHERE id=$7
	`
	_, err := database.DB.Exec(
		ctx, query,
		user.Username, user.Email, user.FullName,
		user.RoleID, user.IsActive, time.Now(),
		user.ID,
	)

	return err
}

// Delete menghapus user (admin delete user).
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	command := `DELETE FROM users WHERE id=$1`
	_, err := database.DB.Exec(ctx, command, id)
	return err
}

// AssignRole mengganti role user (FR-009).
func (r *UserRepository) AssignRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) error {
	query := `UPDATE users SET role_id=$1, updated_at=$2 WHERE id=$3`
	_, err := database.DB.Exec(ctx, query, roleID, time.Now(), id)
	return err
}

// GetAllUsers mengambil semua user untuk admin (FR-009).
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FullName, &user.RoleID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
