package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app *fiber.App) {
	r := app.Group("/api/v1/achievements")

	// protected endpoints
	r.Use(middleware.JWTMiddleware)

	// Mahasiswa: create / update / delete / submit / attachments
	r.Post("/", services.CreateAchievement)                      // Create
	r.Put("/:id", services.UpdateAchievement)                    // Update
	r.Delete("/:id", services.DeleteAchievement)                 // Delete (soft)
	r.Post("/:id/submit", services.SubmitAchievement)            // Submit for verification
	r.Post("/:id/attachments", services.AddAttachment)           // Upload attachments

	// Dosen Wali: verify & reject (routes protected but should also add role-check middleware)
	r.Post("/:id/verify", middleware.RequireRoles("Dosen Wali")(services.VerifyAchievement))
	r.Post("/:id/reject", middleware.RequireRoles("Dosen Wali")(services.RejectAchievement))

	// Common: list / detail / history
	r.Get("/", services.ListAchievements)
	r.Get("/:id", services.GetAchievementDetail)
	r.Get("/:id/history", services.GetAchievementHistory)
}
