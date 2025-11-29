package services

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

// ======================================================
// LOGIN SERVICE
// ======================================================
func Login(ctx context.Context, identifier string, password string) (*models.AuthResponse, error) {

	// 1. Ambil user berdasarkan username/email
	user, err := repositories.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 2. Cek user aktif
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	// 3. Validasi password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid username or password")
	}

	// 4. Ambil Role (pakai .String())
	role, err := repositories.GetRoleByID(ctx, user.RoleID.String())
	if err != nil {
		return nil, errors.New("role not found")
	}

	// 5. Ambil permission_id dari role_permissions
	permissionIDs, err := repositories.GetPermissionIDsByRoleID(ctx, user.RoleID.String())
	if err != nil {
		return nil, err
	}

	// 6. Ambil detail permission
	permissions, err := repositories.GetPermissionsByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}

	// 7. Generate JWT
	tokenString, err := generateJWT(user, role.Name, permissions)
	if err != nil {
		return nil, err
	}

	// 8. Build Response
	return &models.AuthResponse{
		Token:       tokenString,
		User:        user,
		Role:        role,
		Permissions: permissions,
	}, nil
}

// ======================================================
// GENERATE JWT
// ======================================================
func generateJWT(user *models.User, roleName string, perms []models.Permission) (string, error) {

	permStrings := []string{}
	for _, p := range perms {
		permStrings = append(permStrings, p.Name)
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"username":   user.Username,
		"role":       roleName,
		"permissions": permStrings,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
