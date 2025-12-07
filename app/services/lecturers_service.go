package services

import (
	"context"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/v1/lecturers
func GetAllLecturers(c *fiber.Ctx) error {
	ctx := context.Background()

	list, err := repositories.GetAllLecturers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": list})
}

// GET /api/v1/lecturers/:id/advisees
func GetLecturerAdvisees(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	lid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	list, err := repositories.GetAdviseesByLecturerID(ctx, lid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "advisees": list})
}
