package middleware

import "github.com/gofiber/fiber/v2"

// RequirePermission
// FR-002: RBAC Middleware
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		rawPerms := c.Locals("permissions")
		if rawPerms == nil {
			return forbidden(c, "no permissions found")
		}

		perms, ok := rawPerms.([]interface{})
		if !ok {
			return forbidden(c, "invalid permissions format")
		}

		for _, p := range perms {
			if p.(string) == permission {
				return c.Next()
			}
		}

		return forbidden(c, "missing permission: "+permission)
	}
}

func forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": msg,
	})
}
