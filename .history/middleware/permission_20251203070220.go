package middleware

import "github.com/gofiber/fiber/v2"

func RequirePermission(permission string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        raw := c.Locals("permissions")
        if raw == nil {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "forbidden: admin only access",
            })
        }

        perms := raw.([]interface{}) // array of permissions

        for _, p := range perms {
            if p.(string) == permission {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "forbidden: admin only access",
        })
    }
}
