package repository

import (
    "context"
    "database/sql"

    "UAS-backend/models"
)

type UserRepository interface {
    Create(ctx context.Context, u *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    List(ctx context.Context) ([]models.User, error)
    Update(ctx context.Context, u *models.User) error
    Delete(ctx context.Context, id string) error
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db}
}

//
// Create()
// --------------------------------------
// Fungsi ini membuat user baru di database.
// Biasanya dipanggil oleh fitur registrasi / create user oleh admin.
//
func (r *userRepository) Create(ctx context.Context, u *models.User) error {
    query := `
        INSERT INTO users (id, username, email, password_hash, full_name, role_id)
        VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := r.db.ExecContext(ctx, query,
        u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID,
    )
    return err
}

//
// GetByID()
// --------------------------------------
// Mengambil data user berdasarkan ID.
// Digunakan saat menampilkan profil user atau validasi identitas.
//
func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    query := `SELECT id, username, email, full_name, role_id, password_hash, is_active FROM users WHERE id=$1`
    row := r.db.QueryRowContext(ctx, query, id)

    var u models.User
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.PasswordHash, &u.IsActive)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

//
// GetByUsername()
// --------------------------------------
// Mengambil user melalui username.
// Digunakan saat proses login (cek username terdaftar).
//
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
    query := `SELECT id, username, email, full_name, role_id, password_hash, is_active FROM users WHERE username=$1`
    row := r.db.QueryRowContext(ctx, query, username)

    var u models.User
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.PasswordHash, &u.IsActive)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

//
// GetByEmail()
// --------------------------------------
// Mengambil user melalui email.
// Dipakai pada login (mengizinkan login via email) atau fitur cek duplikasi email.
//
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `SELECT id, username, email, full_name, role_id, password_hash, is_active FROM users WHERE email=$1`
    row := r.db.QueryRowContext(ctx, query, email)

    var u models.User
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.PasswordHash, &u.IsActive)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

//
// List()
// --------------------------------------
// Mengambil semua data user.
// Admin menggunakan ini untuk menampilkan daftar semua user.
//
func (r *userRepository) List(ctx context.Context) ([]models.User, error) {
    query := `SELECT id, username, email, full_name, role_id, is_active FROM users`
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive)
        if err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

//
// Update()
// --------------------------------------
// Mengubah data user.
// Umumnya dipakai oleh admin atau fitur edit profile.
//
func (r *userRepository) Update(ctx context.Context, u *models.User) error {
    query := `
        UPDATE users 
        SET username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=NOW()
        WHERE id=$6`
    
    _, err := r.db.ExecContext(ctx, query,
        u.Username, u.Email, u.FullName, u.RoleID, u.IsActive, u.ID,
    )
    return err
}

//
// Delete()
// --------------------------------------
// Menghapus user berdasarkan ID.
// Admin yang menjalankan proses delete user.
//
func (r *userRepository) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM users WHERE id=$1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
