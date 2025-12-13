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

	r.Use(middleware.JWTMiddleware)
	r.Post("/logout", services.LogoutHandler)
	r.Get("/profile", services.GetProfileHandler)
}

