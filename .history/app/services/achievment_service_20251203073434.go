package services

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
)

// ====================================================================
// POST /api/v1/achievements (Mahasiswa only) — Create Draft
// ====================================================================
func CreateAchievement(c *fiber.Ctx) error {
	var req models.CreateAchievementDTO

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()
	studentID := c.Locals("userID").(string)

	req.StudentID = studentID

	// 1. Create in MongoDB
	mongoID, err := repositories.CreateAchievementMongo(ctx, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Create reference in PostgreSQL
	err = repositories.CreateAchievementReference(ctx, mongoID, studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"id":      mongoID,
		"message": "achievement draft created",
	})
}

// ====================================================================
// POST /api/v1/achievements/:id/submit — Submit Achievement
// ====================================================================
func SubmitAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("userID").(string)
	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.StudentID != studentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	err = repositories.UpdateAchievementStatus(ctx, id, "submitted")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement submitted",
	})
}

// ====================================================================
// PUT /api/v1/achievements/:id — Update Draft
// ====================================================================
func UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("userID").(string)

	var req models.UpdateAchievementDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.StudentID != studentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	err = repositories.UpdateAchievementMongo(ctx, id, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement updated",
	})
}

// ====================================================================
// DELETE /api/v1/achievements/:id — Delete
// ====================================================================
func DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	studentID := c.Locals("userID").(string)
	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.StudentID != studentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	err = repositories.SoftDeleteAchievementReference(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err = repositories.SoftDeleteAchievementMongo(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement deleted",
	})
}

// ====================================================================
// POST /api/v1/achievements/:id/verify — Verify (Dosen)
// ====================================================================
func VerifyAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be verified"})
	}

	err = repositories.UpdateAchievementStatus(ctx, id, "verified")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement verified",
	})
}

// ====================================================================
// POST /api/v1/achievements/:id/reject — Reject (Dosen)
// ====================================================================
func RejectAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Note string `json:"note"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
	}

	err = repositories.RejectAchievement(ctx, id, req.Note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement rejected",
	})
}

// ====================================================================
// GET /api/v1/achievements/:id — Detail
// ====================================================================
func GetAchievementDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	ref, err := repositories.GetAchievementReference(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	detail, err := repositories.GetAchievementMongo(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":     id,
			"status": ref.Status,
			"note":   ref.RejectionNote,
			"detail": detail,
		},
	})
}

// ====================================================================
// GET /api/v1/achievements — List (Mahasiswa only OR Dosen/Admin)
// ====================================================================
func ListAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	var result interface{}
	var err error

	if role == "Mahasiswa" {
		result, err = repositories.ListByStudent(ctx, userID)
	} else {
		result, err = repositories.ListAllAchievements(ctx)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
