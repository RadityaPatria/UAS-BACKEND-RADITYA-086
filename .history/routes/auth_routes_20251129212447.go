package routes

import (
	"UAS-backend/app/service"
	"UAS-backend/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {

	auth := app.Group("/api/v1/auth")

	// PUBLIC
	auth.Post("/login", service.Login)

	// PROTECTED (perlu JWT)
	protected := auth.Group("/", middleware.JWTMiddleware)
	protected.Get("/profile", service.GetProfile)
	protected.Post("/logout", service.Logout)
}
