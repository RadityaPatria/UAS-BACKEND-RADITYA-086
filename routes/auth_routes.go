package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {
	r := app.Group("/api/v1/auth")

	// POST /login -> login user | FR-001
	r.Post("/login", services.LoginHandler)

	// POST /refresh -> refresh token | FR-001
	r.Post("/refresh", services.RefreshTokenHandler)

	r.Use(middleware.JWTMiddleware)

	// POST /logout -> logout user | FR-001
	r.Post("/logout", services.LogoutHandler)

	// GET /profile -> profil user login | FR-001
	r.Get("/profile", services.GetProfileHandler)
}
