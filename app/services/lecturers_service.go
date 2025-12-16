package services

import (
	"context"

	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary      Get all lecturers
// @Description  Ambil semua data dosen
// @Tags         Lecturers
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /lecturers [get]
//
// GetAllLecturers -> ambil semua dosen | FR-013
func GetAllLecturers(c *fiber.Ctx) error {
	ctx := context.Background()

	list, err := repositories.GetAllLecturers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}

// @Summary      Get lecturer advisees
// @Description  Ambil daftar mahasiswa bimbingan dosen
// @Tags         Lecturers
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Lecturer ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lecturers/{id}/advisees [get]
//
// GetLecturerAdvisees -> ambil mahasiswa bimbingan dosen | FR-014
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

	return c.JSON(fiber.Map{
		"status":   "success",
		"advisees": list,
	})
}
