package services

import (
	"context"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


// @Summary      Achievement Statistics
// @Description  Statistik pencapaian prestasi (total, status, dsb)
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /reports/achievements/statistics [get]
//
// GetAchievementStatistics -> statistik pencapaian prestasi | FR-015
func GetAchievementStatistics(c *fiber.Ctx) error {
	ctx := context.Background()

	stats, err := repositories.GetAchievementStatistics(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

// @Summary      Student Achievement Report
// @Description  Laporan prestasi per mahasiswa
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Student ID (UUID)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /reports/students/{id} [get]
//
// GetStudentReport -> laporan prestasi per mahasiswa | FR-016
func GetStudentReport(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	sid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	data, err := repositories.GetStudentReport(ctx, sid.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}
