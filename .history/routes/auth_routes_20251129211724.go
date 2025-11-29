package routes

import (
	"context"

	"UAS-backend/app/services"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router) {

	r.Post("/login", func(c *fiber.Ctx) error {
		var req struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		resp, err := services.Login(context.Background(), req.Identifier, req.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   resp,
		})
	})
}
