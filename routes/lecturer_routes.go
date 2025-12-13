package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterLecturerRoutes(app *fiber.App) {
	r := app.Group("/api/v1/lecturers")

	r.Use(middleware.JWTMiddleware)

	// GET / -> list dosen (Admin) | FR-004
	r.Get("/",
		middleware.RequireRoles("Admin"),
		services.GetAllLecturers)

	// GET /:id/advisees -> mahasiswa bimbingan | FR-006
	r.Get("/:id/advisees",
		middleware.RequireRoles("Admin", "Dosen Wali"),
		services.GetLecturerAdvisees)
}
