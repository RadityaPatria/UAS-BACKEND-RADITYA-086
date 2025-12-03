package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")

func JWTMiddleware(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	c.Locals("userID", claims["user_id"])
	c.Locals("role", claims["role"])
	c.Locals("permissions", claims["permissions"])

	return c.Next()
}
