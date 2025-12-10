package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app *fiber.App) {
	r := app.Group("/api/v1/achievements")

	// 1️⃣ SEMUA HARUS LEWAT JWT DULU
	r.Use(middleware.JWTMiddleware)

	// 2️⃣ MAHASISWA
	r.Post("/", services.CreateAchievement)
	r.Put("/:id", services.UpdateAchievement)
	r.Delete("/:id", services.DeleteAchievement)
	r.Post("/:id/submit", services.SubmitAchievement)
	r.Post("/:id/attachments", services.AddAttachment)

	// 3️⃣ DOSEN WALI (HARUS JWT → RequireRoles berurutan)
	r.Post("/:id/verify",
		middleware.RequireRoles("Dosen Wali"),
		services.VerifyAchievement)

	r.Post("/:id/reject",
		middleware.RequireRoles("Dosen Wali"),
		services.RejectAchievement)

	// 4️⃣ COMMON
	r.Get("/", services.ListAchievements)
	r.Get("/:id", services.GetAchievementDetail)
	r.Get("/:id/history", services.GetAchievementHistory)
}
