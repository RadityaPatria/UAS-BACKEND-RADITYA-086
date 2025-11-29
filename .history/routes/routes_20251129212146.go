package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {

	RegisterAuthRoutes(app)
	RegisterUserRoutes(app)
	RegisterAchievementRoutes(app)

}
