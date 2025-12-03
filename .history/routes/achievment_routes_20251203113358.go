package routes

import (
    "UAS-backend/app/services"
    "UAS-backend/app/middleware"

    "github.com/gofiber/fiber/v2"
)

func RegisterAchievementRoutes(app *fiber.App) {

    r := app.Group("/api/v1/achievements")

    // =============================
    // 1. Create Achievement (Mahasiswa)
    // =============================
    r.Post("/", 
        middleware.JWTMiddleware, 
        middleware.RequireRoles("Mahasiswa"),
        services.CreateAchievementHandler,
    )

    // =============================
    // 2. Update Achievement (Draft only)
    // =============================
    r.Put("/:id",
        middleware.JWTMiddleware,
        middleware.RequireRoles("Mahasiswa"),
        services.UpdateAchievementHandler,
    )

    // =============================
    // 3. Submit Achievement
    // =============================
    r.Post("/:id/submit",
        middleware.JWTMiddleware,
        middleware.RequireRoles("Mahasiswa"),
        services.SubmitAchievementHandler,
    )

    // =============================
    // 4. Verify (Dosen Wali)
    // =============================
    r.Post("/:id/verify",
        middleware.JWTMiddleware,
        middleware.RequireRoles("Dosen Wali"),
        services.VerifyAchievementHandler,
    )

    // =============================
    // 5. Reject (Dosen Wali)
    // =============================
    r.Post("/:id/reject",
        middleware.JWTMiddleware,
        middleware.RequireRoles("Dosen Wali"),
        services.RejectAchievementHandler,
    )

    // =============================
    // 6. Delete (Draft only, Mahasiswa)
    // =============================
    r.Delete("/:id",
        middleware.JWTMiddleware,
        middleware.RequireRoles("Mahasiswa"),
        services.DeleteAchievementHandler,
    )

    // =============================
    // 7. Detail Achievement (All Roles)
    // =============================
    r.Get("/:id",
        middleware.JWTMiddleware,
        services.GetAchievementDetailHandler,
    )

    // =============================
    // 8. List Achievements (Role-based)
    // Admin → all, Dosen → advisee only, Mahasiswa → own only
    // =============================
    r.Get("/",
        middleware.JWTMiddleware,
        services.ListAchievementsHandler,
    )
}
