package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	r := app.Group("/api/v1/students")
	r.Use(middleware.JWTMiddleware)

	r.Get("/", services.GetAllStudents)
	r.Get("/:id", services.GetStudentByID)
	r.Get("/:id/achievements", services.GetStudentAchievements)
	r.Put("/:id/advisor", middleware.RequireRoles("Admin"), services.UpdateStudentAdvisor)
}
