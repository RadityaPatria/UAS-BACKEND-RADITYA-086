package routes

import (
    "UAS-backend/app/services"
    "UAS-backend/middleware"

    "github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
    r := app.Group("/api/v1/users")

    // Hanya Admin yang boleh akses seluruh CRUD user
    user.Use(middleware.JWTMiddleware)
    r.Use(middleware.RequireRoles("Admin"))

    r.Get("/", services.GetAllUsers)
    r.Get("/:id", services.GetUserByID)
    r.Post("/", services.CreateUser)
    r.Put("/:id", services.UpdateUser)
    r.Delete("/:id", services.DeleteUser)
    r.Put("/:id/role", services.UpdateUserRole)
}
