package repositor

import (
	"context"
	"UAS-backend/database"
	"UAS-backend/app/models"
	"github.com/google/uuid"
)


// ===============================
// CreateUser()
// Menambahkan user baru.
// Dipakai oleh admin atau fitur registrasi.
// ===============================
func CreateUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := database.DB.Exec(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName,
		u.RoleID, u.IsActive,
	)

	return err
}


// ===============================
// GetUserByID()
// Mengambil user berdasarkan ID.
// Dipakai saat mengambil profil user.
// ===============================
func GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
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


// ===============================
// GetUserByEmail()
// Digunakan pada proses login (validasi email).
// ===============================
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active,
		       created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := database.DB.QueryRow(ctx, query, email)

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


// ===============================
// GetAllUsers()
// Mengambil semua user.
// Dipakai oleh admin.
// ===============================
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active,
		       created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash,
			&u.FullName, &u.RoleID, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}


// ===============================
// UpdateUser()
// Memperbarui data user (profile atau role).
// ===============================
func UpdateUser(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users
		SET username = $2,
		    email = $3,
		    full_name = $4,
		    role_id = $5,
		    is_active = $6,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.DB.Exec(ctx, query,
		u.ID, u.Username, u.Email, u.FullName, u.RoleID, u.IsActive,
	)

	return err
}


// ===============================
// DeleteUser()
// Menghapus user (hard delete).
// Dipakai oleh admin.
// ===============================
func DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, id)
	return err
}
