package services

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

// ==============================
// HTTP HANDLER: LOGIN
// ==============================
func LoginHandler(c *fiber.Ctx) error {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	ctx := context.Background()

	resp, err := Login(ctx, input.Identifier, input.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// ==============================
// LOGIN SERVICE (business logic)
// ==============================
func Login(ctx context.Context, identifier string, password string) (*models.AuthResponse, error) {

	// 1. Ambil user
	user, err := repositories.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 2. Cek aktif
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	// 3. Cek password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid username or password")
	}

	// 4. Ambil role
	role, err := repositories.GetRoleByID(ctx, user.RoleID.String())
	if err != nil {
		return nil, errors.New("role not found")
	}

	// 5. Permission ID
	permissionIDs, err := repositories.GetPermissionIDsByRoleID(ctx, user.RoleID.String())
	if err != nil {
		return nil, err
	}

	// 6. Detail permission
	permissions, err := repositories.GetPermissionsByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, err
	}

	// 7. Generate JWT
	tokenString, err := generateJWT(user, role.Name, permissions)
	if err != nil {
		return nil, err
	}

	// 8. Return
	return &models.AuthResponse{
		Token:       tokenString,
		User:        user,
		Role:        role,
		Permissions: permissions,
	}, nil
}

// ==============================
// JWT Maker
// ==============================
func generateJWT(user *models.User, roleName string, perms []models.Permission) (string, error) {

	var permStrings []string
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
