package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Running 🚀")
	})
}
package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {

	RegisterAuthRoutes(app)
	RegisterUserRoutes(app)
	RegisterAchievementRoutes(app)

}
