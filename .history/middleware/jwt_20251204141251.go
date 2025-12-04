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

	claims := token.Claims.(jwt.MapClaims)

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	c.Locals("userID", userID)
	c.Locals("role", role)

	// Map student
	if role == "Mahasiswa" {
		studentID, _ := repositories.GetStudentIDByUserID(context.Background(), userID)
		if studentID != "" {
			c.Locals("studentID", studentID)
		}
	}

	// Map lecturer (IMPORTANT: use lecturers.id, not userID!)
	if role == "Dosen Wali" {
		lecturerID, _ := repositories.GetLecturerIDByUserID(context.Background(), userID)
		if lecturerID != "" {
			c.Locals("lecturerID", lecturerID)
		}
	}

	return c.Next()
}
