package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		perms := c.Locals("permissions")
		if perms == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permissions not found",
			})
		}

		list := perms.([]interface{})
		for _, p := range list {
			if p.(string) == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden: missing permission " + permission,
		})
	}
}
