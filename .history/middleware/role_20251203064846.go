package middlewares

import "github.com/gofiber/fiber/v2"

// RequireRole("admin")
// RequireRole("mahasiswa", "dosen_wali")
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		currentRole := c.Locals("role")
		if currentRole == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "role not found in token",
			})
		}

		roleStr := currentRole.(string)

		for _, allowed := range roles {
			if roleStr == allowed {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "forbidden: role not allowed",
		})
	}
}
