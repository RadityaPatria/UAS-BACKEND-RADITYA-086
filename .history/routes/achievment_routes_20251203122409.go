package routes

import (
    "UAS-backend/app/services"
    "UAS-backend/middleware"

    "github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app *fiber.App) {

    r := app.Group("/api/v1/achievements")

    // All routes require login
    r.Use(middleware.JWTMiddleware)

    // ============================================================
    // Mahasiswa only — CREATE, UPDATE, DELETE, SUBMIT
    // ============================================================
    r.Post("/", middleware.RequireRoles("Mahasiswa"), services.CreateAchievement)
    r.Put("/:id", middleware.RequireRoles("Mahasiswa"), services.UpdateAchievement)
    r.Delete("/:id", middleware.RequireRoles("Mahasiswa"), services.DeleteAchievement)
    r.Post("/:id/submit", middleware.RequireRoles("Mahasiswa"), services.SubmitAchievement)

    // ============================================================
    // Dosen Wali — VERIFY & REJECT
    // ============================================================
    r.Post("/:id/verify", middleware.RequireRoles("Dosen Wali"), services.VerifyAchievement)
    r.Post("/:id/reject", middleware.RequireRoles("Dosen Wali"), services.RejectAchievement)

    // ============================================================
    // Admin, Dosen Wali, Mahasiswa — LIST & DETAIL
    // ============================================================
    r.Get("/", middleware.RequireRoles("Admin", "Dosen Wali", "Mahasiswa"), services.ListAchievements)
    r.Get("/:id", middleware.RequireRoles("Admin", "Dosen Wali", "Mahasiswa"), services.GetAchievementDetail)

    // ============================================================
    // Riwayat
    // ============================================================
    r.Get("/:id/history", middleware.RequireRoles("Admin", "Dosen Wali", "Mahasiswa"), services.GetAchievementHistory)

    // ============================================================
    // Upload attachment (Mahasiswa only)
    // ============================================================
    r.Post("/:id/attachments", middleware.RequireRoles("Mahasiswa"), services.AddAttachment)
}
