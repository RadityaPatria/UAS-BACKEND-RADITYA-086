package middlewares

import "github.com/gofiber/fiber/v2"

// RequirePermission("user:manage")
// RequirePermission("achievement:create")
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		raw := c.Locals("permissions")
		if raw == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "permissions not found in token",
			})
		}

		perms := raw.([]interface{})

		for _, p := range perms {
			if p.(string) == permission {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "forbidden: missing permission (" + permission + ")",
		})
	}
}
