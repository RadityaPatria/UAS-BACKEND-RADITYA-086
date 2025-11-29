package repositories

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)


// GetUserByIdentifier → ambil user pakai username atau email
func GetUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		 FROM users WHERE username=$1 OR email=$1`, identifier)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}


// GetUserByID → ambil user berdasarkan ID
func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		 FROM users WHERE id=$1`, id)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}


// CreateUser → buat user baru (ADMIN)
func CreateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash, full_name, role_id, 
			is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())`,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.FullName, u.RoleID, u.IsActive,
	)
	return err
}

//
// UpdateUser → update data user (ADMIN)

func UpdateUser(ctx context.Context, u *models.User) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE users SET username=$1, email=$2, full_name=$3, role_id=$4, 
		 is_active=$5, updated_at=NOW() WHERE id=$6`,
		u.Username, u.Email, u.FullName, u.RoleID,
		u.IsActive, u.ID,
	)
	return err
}

//
// DeleteUser → hapus user (ADMIN)
//
func DeleteUser(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
