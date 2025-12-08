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
	c.Locals("permissions", claims["permissions"])

	// MAP STUDENT ID
	if role == "Mahasiswa" && userID != "" {
		sid, err := repositories.GetStudentIDByUserID(context.Background(), userID)
		if err == nil && sid != "" {
			c.Locals("studentID", sid)
		}
	}

	// MAP LECTURER ID
	if role == "Dosen Wali" && userID != "" {
		lecturer, err := repositories.GetLecturerByUserID(context.Background(), userID)
		if err != nil || lecturer == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "lecturer profile not found for this user",
			})
		}

		// lecturers.id → UUID
		c.Locals("lecturerID", lecturer.ID.String())
	}

	return c.Next()
}
