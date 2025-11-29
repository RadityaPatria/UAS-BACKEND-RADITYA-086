package routes

import (
	"UAS-backend/app/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	r := app.Group("/api/v1/users")

	r.Get("/", services.GetAllUsers)
	r.Get("/:id", services.GetUserByID)
	r.Post("/", services.CreateUser)
	r.Put("/:id", services.UpdateUser)
	r.Delete("/:id", services.DeleteUser)
	r.Put("/:id/role", services.UpdateUserRole)
}
