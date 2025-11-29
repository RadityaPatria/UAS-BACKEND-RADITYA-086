package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"
	"UAS-backend/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

var ctx = context.Background()

// ============================================================
// LOGIN SERVICE — Validasi user & buat token JWT
// ============================================================

func LoginService(db *sql.DB, identifier string, password string, cfg *config.Config) (*models.LoginResponse, error) {

	// 1. Ambil user via email atau username
	user, err := repositories.GetUserByEmailOrUsername(db, identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Cek apakah user aktif
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	// 3. Cek password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 4. Ambil role name
	roleName, err := repositories.GetRoleNameByID(db, user.RoleID.String())
	if err != nil {
		return nil, errors.New("role not found")
	}

	// 5. Ambil permissions user
	permissions, err := repositories.GetUserPermissions(db, user.ID.String())
	if err != nil {
		return nil, errors.New("failed to fetch permissions")
	}

	// 6. Generate access token
	token, err := generateAccessToken(user.ID, roleName, permissions, cfg)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 7. Generate refresh token
	refreshToken, err := generateRefreshToken(user.ID, cfg)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// 8. Susun response
	resp := &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: models.LoginUser{
			ID:         user.ID,
			Username:   user.Username,
			FullName:   user.FullName,
			Role:       roleName,
			Permissions: permissions,
		},
	}

	return resp, nil
}

//
// ============================================================
// JWT GENERATOR
// ============================================================
//

// TOKEN UTAMA (akses API)
func generateAccessToken(userID uuid.UUID, role string, perms []string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"perms": perms,
		"exp":  time.Now().Add(time.Hour * 1).Unix(), // 1 jam
		"iat":  time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.JWTSecret))
}

// REFRESH TOKEN
func generateRefreshToken(userID uuid.UUID, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 hari
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.JWTRefreshSecret))
}
