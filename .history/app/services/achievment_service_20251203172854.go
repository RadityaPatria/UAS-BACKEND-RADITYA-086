package services

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /api/v1/achievements
// Create achievement (Mahasiswa) -> save to Mongo + reference to Postgres
func CreateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()

	// determine studentID from Locals (set by middleware)
	studentIDraw := c.Locals("studentID")
	userIDraw := c.Locals("userID")

	var studentID string
	if studentIDraw != nil {
		studentID = studentIDraw.(string)
	} else if userIDraw != nil {
		// fallback: try mapping
		userID := userIDraw.(string)
		sid, err := repositories.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "student record not found for this user"})
		}
		studentID = sid
	} else {
		return c.Status(401).JSON(fiber.Map{"error": "unauthenticated"})
	}

	var payload models.AchievementMongo
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// ensure basic fields
	payload.Status = models.StatusDraft
	payload.StudentID = studentID

	// save to Mongo
	mongoID, err := repositories.CreateAchievementMongo(ctx, &payload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// create postgres reference
	ref := models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.MustParse(studentID),
		MongoAchievementID: mongoID.Hex(),
		Status:             models.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := repositories.CreateAchievementReference(ctx, &ref); err != nil {
		// rollback: try delete mongo (best effort)
		_ = repositories.SoftDeleteAchievementMongo(ctx, mongoID)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// build return achievement: combine mongo + reference (returning reference is OK per spec)
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"reference": ref,
			"mongo_id":  mongoID.Hex(),
		},
	})
}

// PUT /api/v1/achievements/:id
// Update achievement in Mongo (Mahasiswa, only DRAFT)
func UpdateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	// parse body as map (bson.M)
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	// check ownership
	studentLocal := c.Locals("studentID")
	if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	// convert map -> bson.M
	updateBson := bson.M(update)
	if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateBson); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// update timestamp on reference
	if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusDraft); err != nil {
		// ignore non-fatal
	}

	return c.JSON(fiber.Map{"status": "success", "message": "updated"})
}

// POST /api/v1/achievements/:id/submit
func SubmitAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	// ownership
	studentLocal := c.Locals("studentID")
	if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	// set submitted
	if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusSubmitted); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// set submitted_at in db (GetAchievementReferenceByID returns updated when fetched later)
	return c.JSON(fiber.Map{"status": "success", "message": "submitted"})
}

// DELETE /api/v1/achievements/:id
func DeleteAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	studentLocal := c.Locals("studentID")
	if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	if err := repositories.SoftDeleteAchievementReference(ctx, refID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err := repositories.SoftDeleteAchievementMongo(ctx, mongoID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "deleted"})
}

// GET /api/v1/achievements
func ListAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	var list []models.AchievementReference
	var err error

	switch role {
	case "Admin":
		list, err = repositories.ListAllAchievementReferences(ctx)
	case "Dosen Wali":
		lecturerUUID, err2 := uuid.Parse(userID)
		if err2 != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid lecturer id"})
		}
		advisees, err2 := repositories.GetAdviseesByLecturerID(ctx, lecturerUUID)
		if err2 != nil {
			return c.Status(500).JSON(fiber.Map{"error": err2.Error()})
		}
		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)
	case "Mahasiswa":
		studentLocal := c.Locals("studentID")
		if studentLocal == nil {
			return c.Status(400).JSON(fiber.Map{"error": "student mapping not found"})
		}
		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, []string{studentLocal.(string)})
	default:
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": list})
}

// POST /api/v1/achievements/:id/verify  (Dosen Wali)
func VerifyAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	lecturerID := c.Locals("userID").(string)

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if ref.Status != models.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be verified"})
	}

	if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := repositories.SetVerifiedBy(ctx, refID, lecturerID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "verified"})
}

// POST /api/v1/achievements/:id/reject  (Dosen Wali)
func RejectAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if ref.Status != models.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
	}
	if err := repositories.SetAchievementRejectionNote(ctx, refID, body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "rejected"})
}

// GET /api/v1/achievements/:id
func GetAchievementDetail(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	ach, err := repositories.GetAchievementMongoByID(ctx, mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "reference": ref, "achievement": ach})
}

// GET /api/v1/achievements/:id/history
func GetAchievementHistory(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "history": ref})
}

// POST /api/v1/achievements/:id/attachments
func AddAttachment(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	var att models.Attachment
	if err := c.BodyParser(&att); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err := repositories.AddAttachmentToAchievement(ctx, mongoID, att); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "attachment added"})
}
