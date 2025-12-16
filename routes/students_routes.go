package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	r := app.Group("/api/v1/students")

	r.Use(middleware.JWTMiddleware)

	// GET / -> list mahasiswa | FR-004
	r.Get("/", services.GetAllStudents)

	// GET /:id -> detail mahasiswa | FR-004
	r.Get("/:id", services.GetStudentByID)

	// GET /:id/achievements -> prestasi mahasiswa | FR-007
	r.Get("/:id/achievements", services.GetStudentAchievements)

	// PUT /:id/advisor -> update dosen wali (Admin) | FR-004
	r.Put("/:id/advisor",
		middleware.RequireRoles("Admin"), services.UpdateStudentAdvisor)
}
