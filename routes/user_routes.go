package routes

import (
	"UAS-backend/app/services"
	"UAS-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	r := app.Group("/api/v1/users")

	r.Use(middleware.JWTMiddleware)
	r.Use(middleware.RequireRoles("Admin"))

	// GET / -> list semua user | FR-002
	r.Get("/", services.GetAllUsers)

	// GET /:id -> detail user | FR-002
	r.Get("/:id", services.GetUserByID)

	// POST / -> tambah user | FR-002
	r.Post("/", services.CreateUser)

	// PUT /:id -> update user | FR-002
	r.Put("/:id", services.UpdateUser)

	// DELETE /:id -> hapus user | FR-002
	r.Delete("/:id", services.DeleteUser)

	// PUT /:id/role -> ubah role user | FR-002
	r.Put("/:id/role", services.UpdateUserRole)
}
