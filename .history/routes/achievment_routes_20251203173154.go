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
	r.Post("/", services.CreateAchievement)                      
	r.Put("/:id", services.UpdateAchievement)                    
	r.Delete("/:id", services.DeleteAchievement)                 
	r.Post("/:id/submit", services.SubmitAchievement)            
	r.Post("/:id/attachments", services.AddAttachment)           

	// Dosen Wali: verify & reject
	r.Post("/:id/verify", middleware.RequireRoles("Dosen Wali"), services.VerifyAchievement)
	r.Post("/:id/reject", middleware.RequireRoles("Dosen Wali"), services.RejectAchievement)

	// Common: list / detail / history
	r.Get("/", services.ListAchievements)
	r.Get("/:id", services.GetAchievementDetail)
	r.Get("/:id/history", services.GetAchievementHistory)
}
