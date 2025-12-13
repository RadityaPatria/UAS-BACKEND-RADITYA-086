package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterLecturerRoutes(app *fiber.App) {
    r := app.Group("/api/v1/lecturers")
    r.Use(middleware.JWTMiddleware)

    // Admin only
    r.Get("/", middleware.RequireRoles("Admin"), services.GetAllLecturers)

    // Admin & dosen wali
    r.Get("/:id/advisees",
        middleware.RequireRoles("Admin", "Dosen Wali"),
        services.GetLecturerAdvisees)
}

