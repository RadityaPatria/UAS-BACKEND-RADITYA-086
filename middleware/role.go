package middleware

import "github.com/gofiber/fiber/v2"

// rolesAllowed = ["Admin", "Dosen", "Mahasiswa"]
func RequireRoles(rolesAllowed ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        rawRole := c.Locals("role")
        if rawRole == nil {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "forbidden: role not found",
            })
        }

        userRole := rawRole.(string)

        for _, allowed := range rolesAllowed {
            if userRole == allowed {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "forbidden: only " + joinRoles(rolesAllowed) + " can access this endpoint",
        })
    }
}

func joinRoles(roles []string) string {
    if len(roles) == 1 {
        return roles[0]
    }
    res := roles[0]
    for i := 1; i < len(roles); i++ {
        res += ", " + roles[i]
    }
    return res
}
