package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireRoles
// FR-002: RBAC Middleware
func RequireRoles(rolesAllowed ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		rawRole := c.Locals("role")
		if rawRole == nil {
			return forbidden(c, "role not found")
		}

		userRole := rawRole.(string)

		for _, role := range rolesAllowed {
			if userRole == role {
				return c.Next()
			}
		}

		return forbidden(c, "only ["+strings.Join(rolesAllowed, ", ")+"] can access this endpoint")
	}
}
