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

// ======================================================
// LOGIN HANDLER
// ======================================================
func LoginHandler(c *fiber.Ctx) error {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ctx := context.Background()

	resp, err := Login(ctx, req.Identifier, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// ======================================================
// MAIN LOGIN LOGIC
// ======================================================
func Login(ctx context.Context, identifier string, password string) (interface{}, error) {

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

	// 4. Ambil Role
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

	// ======================================================
	// 8. Build Clean Response (HILANGKAN DATA TIDAK PERLU)
	// ======================================================

	// Ambil permission name saja (string)
	cleanPerms := []string{}
	for _, p := range permissions {
		cleanPerms = append(cleanPerms, p.Name)
	}

	// Build user yang sudah bersih
	cleanUser := fiber.Map{
		"id":          user.ID.String(),
		"username":    user.Username,
		"fullName":    user.FullName,
		"role":        role.Name,
		"permissions": cleanPerms,
	}

	// Final response
	return fiber.Map{
		"token": tokenString,
		"user":  cleanUser,
	}, nil
}

// ======================================================
// GET PROFILE HANDLER
// ======================================================
func GetProfileHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID")
	role := c.Locals("role")
	perms := c.Locals("permissions")

	return c.JSON(fiber.Map{
		"user_id":     userID,
		"role":        role,
		"permissions": perms,
	})
}

// ======================================================
// LOGOUT HANDLER
// ======================================================
func LogoutHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "logged out",
	})
}

// ======================================================
// REFRESH TOKEN
// ======================================================
func RefreshTokenHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"token": "refresh-token-not-implemented",
	})
}

// ======================================================
// JWT GENERATOR
// ======================================================
func generateJWT(user *models.User, roleName string, perms []models.Permission) (string, error) {

	permStrings := []string{}
	for _, p := range perms {
		permStrings = append(permStrings, p.Name)
	}

	claims := jwt.MapClaims{
		"user_id":     user.ID.String(),
		"username":    user.Username,
		"role":        roleName,
		"permissions": permStrings,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
