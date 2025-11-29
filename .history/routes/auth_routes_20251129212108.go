package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {

	auth := app.Group("/api/v1/auth")

	// PUBLIC (tanpa token)
	auth.Post("/login", services.LoginHandler)
	auth.Post("/refresh", services.RefreshTokenHandler)

	// PROTECTED (harus login)
	protected := auth.Group("/", middlewares.JWTMiddleware)

	protected.Get("/profile", services.GetProfileHandler)
	protected.Post("/logout", services.LogoutHandler)
}
