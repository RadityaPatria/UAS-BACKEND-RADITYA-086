package repositor

import (
	"context"
	"UAS-backend/app/models"
	"UAS-backend/database"
)

func GetUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, 
		is_active, created_at, updated_at
		FROM users
		WHERE username=$1 OR email=$1
	`

	row := database.DB.QueryRow(ctx, query, identifier)

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

func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, 
		is_active, created_at, updated_at
		FROM users WHERE id=$1
	`

	row := database.DB.QueryRow(ctx, query, id)

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

func CreateUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, 
			is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
	`

	_, err := database.DB.Exec(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.FullName, u.RoleID, u.IsActive,
	)

	return err
}

func UpdateUser(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users SET username=$1, email=$2, full_name=$3, role_id=$4,
		is_active=$5, updated_at=NOW()
		WHERE id=$6
	`

	_, err := database.DB.Exec(ctx, query,
		u.Username, u.Email, u.FullName, u.RoleID,
		u.IsActive, u.ID,
	)

	return err
}

func DeleteUser(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
