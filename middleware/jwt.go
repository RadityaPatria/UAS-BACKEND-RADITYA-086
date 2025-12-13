package middleware

import (
	"context"
	"strings"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT secret (sebaiknya dari env di production)
var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

// JWTMiddleware
// FR-002: RBAC Middleware
func JWTMiddleware(c *fiber.Ctx) error {
	// ===== Extract Authorization Header =====
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return unauthorized(c, "missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// ===== Parse & Validate Token =====
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return unauthorized(c, "invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return unauthorized(c, "invalid token claims")
	}

	// ===== Extract Claims =====
	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	tokenVersionFloat, ok := claims["token_version"].(float64)
	if !ok {
		return unauthorized(c, "invalid token version")
	}
	tokenVersion := int(tokenVersionFloat)

	// ===== Token Revocation Check =====
	// FR-001: Logout / Refresh Security
	user, err := repositories.GetUserByID(context.Background(), userID)
	if err != nil || user.TokenVersion != tokenVersion {
		return unauthorized(c, "token has been revoked")
	}

	// ===== Inject Context =====
	c.Locals("userID", userID)
	c.Locals("role", role)
	c.Locals("permissions", claims["permissions"])

	// ===== Role Mapping =====
	if role == "Mahasiswa" {
		mapStudent(c, userID)
	}

	if role == "Dosen Wali" {
		if err := mapLecturer(c, userID); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Next()
}

// ===== Helpers =====

func unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": msg,
	})
}

// FR-003, FR-004, FR-005 (Mahasiswa Context)
func mapStudent(c *fiber.Ctx, userID string) {
	sid, err := repositories.GetStudentIDByUserID(context.Background(), userID)
	if err == nil && sid != "" {
		c.Locals("studentID", sid)
	}
}

// FR-006, FR-007, FR-008 (Dosen Wali Context)
func mapLecturer(c *fiber.Ctx, userID string) error {
	lecturer, err := repositories.GetLecturerByUserID(context.Background(), userID)
	if err != nil || lecturer == nil {
		return fiber.NewError(fiber.StatusForbidden, "lecturer profile not found")
	}
	c.Locals("lecturerID", lecturer.ID.String())
	return nil
}
