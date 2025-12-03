package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		current := c.Locals("role")
		if current == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "role not found in token",
			})
		}

		userRole := current.(string)

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": fmt.Sprintf(
				"forbidden: role '%s' tidak diperbolehkan, hanya role [%s] yang bisa akses",
				userRole,
				formatRoles(allowedRoles),
			),
		})
	}
}

func formatRoles(roles []string) string {
	result := ""
	for i, r := range roles {
		if i > 0 {
			result += ", "
		}
		result += r
	}
	return result
}
