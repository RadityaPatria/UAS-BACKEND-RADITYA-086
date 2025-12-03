package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Locals("role")
		if raw == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "role not found in token",
			})
		}

		userRole := raw.(string)

		// cek apakah role diizinkan
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		// jika bukan role yg benar → tolak
		return c.Status(403).JSON(fiber.Map{
			"error": fmt.Sprintf(
				"forbidden: role '%s' tidak diperbolehkan, hanya [%s] yang bisa akses",
				userRole,
				joinRoles(allowedRoles),
			),
		})
	}
}

func joinRoles(roles []string) string {
	out := ""
	for i, r := range roles {
		if i > 0 {
			out += ", "
		}
		out += r
	}
	return out
}
