package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterLecturerRoutes(app *fiber.App) {
	r := app.Group("/api/v1/lecturers")
	r.Use(middleware.JWTMiddleware)

	r.Get("/", middleware.RequireRoles("Admin"), services.GetAllLecturers)
	r.Get("/:id/advisees", middleware.RequireRoles("Admin", "Dosen Wali"), services.GetLecturerAdvisees)
}
