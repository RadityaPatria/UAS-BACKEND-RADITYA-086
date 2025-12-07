package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterReportRoutes(app *fiber.App) {
	r := app.Group("/api/v1/reports")
	r.Use(middleware.JWTMiddleware)

	r.Get("/statistics", middleware.RequireRoles("Admin"), services.GetAchievementStatistics)
	r.Get("/student/:id", middleware.RequireRoles("Admin", "Dosen Wali"), services.GetStudentReport)
}
