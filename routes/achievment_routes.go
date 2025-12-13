package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app *fiber.App) {
	r := app.Group("/api/v1/achievements")

	r.Use(middleware.JWTMiddleware)

	// POST / -> buat prestasi | FR-007
	r.Post("/", services.CreateAchievement)

	// PUT /:id -> update prestasi | FR-007
	r.Put("/:id", services.UpdateAchievement)

	// DELETE /:id -> hapus prestasi | FR-007
	r.Delete("/:id", services.DeleteAchievement)

	// POST /:id/submit -> submit prestasi | FR-007
	r.Post("/:id/submit", services.SubmitAchievement)

	// POST /:id/attachments -> upload lampiran | FR-007
	r.Post("/:id/attachments", services.AddAttachment)

	// POST /:id/verify -> verifikasi prestasi | FR-006
	r.Post("/:id/verify",
		middleware.RequireRoles("Dosen Wali"),
		services.VerifyAchievement)

	// POST /:id/reject -> tolak prestasi | FR-006
	r.Post("/:id/reject",
		middleware.RequireRoles("Dosen Wali"),
		services.RejectAchievement)

	// GET / -> list prestasi | FR-007
	r.Get("/", services.ListAchievements)

	// GET /:id -> detail prestasi | FR-007
	r.Get("/:id", services.GetAchievementDetail)

	// GET /:id/history -> riwayat prestasi | FR-007
	r.Get("/:id/history", services.GetAchievementHistory)
}
