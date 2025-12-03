package middleware

import (
	"context"
	"strings"
	"time"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

// JWTMiddleware parse token, verify, lalu set Locals:
// - userID (users.id as string)
// - role (string)
// - permissions ([]string) raw from claim
// - studentID (string) if role == "Mahasiswa" and mapping exists
// - lecturerID (string) if role == "Dosen Wali"
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	// safe cast user_id & role
	userID := ""
	role := ""
	if v, ok := claims["user_id"].(string); ok {
		userID = v
	}
	if v, ok := claims["role"].(string); ok {
		role = v
	}

	// permissions may be []interface{} in claim
	permRaw := claims["permissions"]

	// set basic locals
	c.Locals("userID", userID)
	c.Locals("role", role)
	c.Locals("permissions", permRaw)

	// if mahasiswa, try map to students.id
	if role == "Mahasiswa" && userID != "" {
		studentID, err := repositories.GetStudentIDByUserID(context.Background(), userID)
		if err == nil && studentID != "" {
			c.Locals("studentID", studentID)
		}
	}

	// if dosen, set lecturerID as userID (most likely lecturers.user_id == users.id)
	if role == "Dosen Wali" && userID != "" {
		c.Locals("lecturerID", userID)
	}

	return c.Next()
}
