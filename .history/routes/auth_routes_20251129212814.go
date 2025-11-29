package routes

import (
	"UAS-backend/services"
	"UAS-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {

	auth := app.Group("/api/v1/auth")

	// PUBLIC
	auth.Post("/login", services.LoginHandler)

	// PROTECTED
	me := auth.Group("/me", middlewares.JWTMiddleware)
	me.Get("/profile", services.GetProfileHandler)
	me.Post("/logout", services.LogoutHandler)
}
