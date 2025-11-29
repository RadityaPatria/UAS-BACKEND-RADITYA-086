package routes

import (
	"UAS-backend/app/services"
	

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {

	r := app.Group("/api/auth")

	// PUBLIC
	r.Post("/login", services.LoginHandler)



}
