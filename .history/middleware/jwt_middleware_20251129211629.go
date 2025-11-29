package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(500).JSON(fiber.Map{"error": "JWT secret not set"})
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	// simpan ke Locals untuk RBAC middleware
	c.Locals("userID", claims["user_id"])
	c.Locals("role", claims["role"])
	c.Locals("permissions", claims["permissions"])

	return c.Next()
}
