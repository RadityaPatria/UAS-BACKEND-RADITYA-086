package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {
	r := app.Group("/api/v1/auth")

	r.Post("/login", services.LoginHandler)
	r.Post("/refresh", services.RefreshTokenHandler)

	// Protected
	r.Post("/logout", middleware.JWTMiddleware, services.LogoutHandler)
	r.Get("/profile", middleware.JWTMiddleware, services.GetProfileHandler)
}
