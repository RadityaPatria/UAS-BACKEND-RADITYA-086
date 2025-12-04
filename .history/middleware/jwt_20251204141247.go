package middleware

import (
	"context"
	"strings"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

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

	userID := ""
	role := ""

	if v, ok := claims["user_id"].(string); ok {
		userID = v
	}
	if v, ok := claims["role"].(string); ok {
		role = v
	}

	c.Locals("userID", userID)
	c.Locals("role", role)

	// Mahasiswa → Map ke students.id
	if role == "Mahasiswa" && userID != "" {
		studentID, _ := repositories.GetStudentIDByUserID(context.Background(), userID)
		if studentID != "" {
			c.Locals("studentID", studentID)
		}
	}

	// Dosen Wali → Map KE lecturer.id (bukan users.id)
	if role == "Dosen Wali" && userID != "" {
		lecturerID, _ := repositories.GetLecturerIDByUserID(context.Background(), userID)
		if lecturerID != "" {
			c.Locals("lecturerID", lecturerID)
		}
	}

	return c.Next()
}
