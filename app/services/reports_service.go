package services

import (
	"context"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GET /api/v1/reports/statistics
func GetAchievementStatistics(c *fiber.Ctx) error {
	ctx := context.Background()

	stats, err := repositories.GetAchievementStatistics(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": stats})
}

// GET /api/v1/reports/student/:id
func GetStudentReport(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	sid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	data, err := repositories.GetStudentReport(ctx, sid.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": data})
}
