package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {

	r := app.Group("/api/auth")

	// PUBLIC
	r.Post("/login", services.LoginHandler)

	// PROTECTED
	me := r.Group("/me", middlewares.JWTMiddleware)
	me.Get("/profile", services.GetProfileHandler)
	me.Get("/logout", services.LogoutHandler)
}
