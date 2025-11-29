package middlewares

import "github.com/gofiber/fiber/v2"

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role := c.Locals("role")

		// Jika admin → lolos
		if role != nil && role.(string) == "Admin" {
			return c.Next()
		}

		// Jika bukan admin → langsung blokir tanpa cek permission
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only admin can access this endpoint",
		})
	}
}
