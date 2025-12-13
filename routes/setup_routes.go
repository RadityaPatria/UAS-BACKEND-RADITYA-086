package routes

import "github.com/gofiber/fiber/v2"

// GET / -> health check API | FR-000
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Running 🚀")
	})
}
