package services

import (
	"context"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ====================================================================
// POST /api/v1/achievements  — Create Draft (Mahasiswa)
// ====================================================================
func CreateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()

	var req struct {
		AchievementType string                 `json:"achievementType"`
		Title           string                 `json:"title"`
		Description     string                 `json:"description"`
		Details         map[string]interface{} `json:"details"`
		Tags            []string               `json:"tags"`
		Points          int                    `json:"points"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// userID from JWT (string UUID)
	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Validate student's UUID (will also be stored in mongo as string)
	if _, err := uuid.Parse(userIDStr); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	// Build Mongo model
	mongoDoc := &models.AchievementMongo{
		StudentID:       userIDStr, // sesuai model: string
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Attachments:     []models.Attachment{},
		Tags:            req.Tags,
		Points:          req.Points,
	}

	// Save to Mongo
	mongoID, err := repositories.CreateAchievementMongo(ctx, mongoDoc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save achievement to mongo"})
	}

	// Create reference in PostgreSQL
	studentUUID, _ := uuid.Parse(userIDStr)
	ref := &models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          studentUUID,
		MongoAchievementID: mongoID.Hex(),
		Status:             models.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := repositories.CreateAchievementReference(ctx, ref); err != nil {
		// attempt to cleanup mongo? optional
		return c.Status(500).JSON(fiber.Map{"error": "failed to save achievement reference"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"referenceId": ref.ID.String(),
			"mongoId":     mongoID.Hex(),
			"status":      ref.Status,
		},
	})
}

// ====================================================================
// PUT /api/v1/achievements/:id — Update Draft (Mahasiswa, only draft)
//    Here :id is reference ID (UUID)
// ====================================================================
func UpdateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	var req struct {
		AchievementType *string                `json:"achievementType"`
		Title           *string                `json:"title"`
		Description     *string                `json:"description"`
		Details         map[string]interface{} `json:"details"`
		Tags            []string               `json:"tags"`
		Points          *int                   `json:"points"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get reference
	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	// Ownership check: studentID in reference vs user
	if ref.StudentID.String() != userIDStr {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	// Must be draft to update
	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	// Convert mongo id to ObjectID
	objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "invalid mongo id stored in reference"})
	}

	// Build update map partial
	update := make(map[string]interface{})
	if req.AchievementType != nil {
		update["achievementType"] = *req.AchievementType
	}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Details != nil {
		update["details"] = req.Details
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}
	if req.Points != nil {
		update["points"] = *req.Points
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "nothing to update"})
	}

	if err := repositories.UpdateAchievementMongo(ctx, objID, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update achievement"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "achievement updated"})
}

// ====================================================================
// DELETE /api/v1/achievements/:id — Soft Delete (Mahasiswa, only draft)
// ====================================================================
func DeleteAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	// Ownership & status check
	if ref.StudentID.String() != userIDStr {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}
	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
	}

	// Soft delete in Postgres
	if err := repositories.SoftDeleteAchievementReference(ctx, refID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to soft delete reference"})
	}

	// Soft delete in Mongo
	objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err == nil {
		_ = repositories.SoftDeleteAchievementMongo(ctx, objID) // ignore mongo error? log optionally
	}

	return c.JSON(fiber.Map{"status": "success", "message": "achievement deleted"})
}

// ====================================================================
// POST /api/v1/achievements/:id/submit — Submit (Mahasiswa, draft -> submitted)
// ====================================================================
func SubmitAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.StudentID.String() != userIDStr {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusSubmitted); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to submit achievement"})
	}

	// TODO: create notification for advisor (async/queue) — not implemented here

	return c.JSON(fiber.Map{"status": "success", "message": "achievement submitted"})
}

// ====================================================================
// POST /api/v1/achievements/:id/verify — Verify (Dosen, submitted -> verified)
// ====================================================================
func VerifyAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	// optional: check role from locals (middleware should already ensure)
	role, _ := c.Locals("role").(string)
	if role != "Dosen" && role != "Admin" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: only Dosen/Admin can verify"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != models.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be verified"})
	}

	// Set verified_by (from token userID) and update status
	verifierIDStr, _ := c.Locals("userID").(string)
	verifierUUID, err := uuid.Parse(verifierIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid verifier id"})
	}

	// update status to verified
	if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to verify achievement"})
	}

	// set verified_by manually via SQL (repository may need this)
	_, _ = repositories.DB.Exec(ctx,
		`UPDATE achievement_references SET verified_by=$2, verified_at=NOW(), updated_at=NOW() WHERE id=$1`,
		refID, verifierUUID,
	) // if you prefer, implement SetVerifiedBy(repo) instead

	return c.JSON(fiber.Map{"status": "success", "message": "achievement verified"})
}

// ====================================================================
// POST /api/v1/achievements/:id/reject — Reject (Dosen, submitted -> rejected)
// ====================================================================
func RejectAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	role, _ := c.Locals("role").(string)
	if role != "Dosen" && role != "Admin" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: only Dosen/Admin can reject"})
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil || req.Note == "" {
		return c.Status(400).JSON(fiber.Map{"error": "rejection note required"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != models.StatusSubmitted {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
	}

	if err := repositories.SetAchievementRejectionNote(ctx, refID, req.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to set rejection note"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "achievement rejected"})
}

// ====================================================================
// POST /api/v1/achievements/:id/attachments — Add Attachment (Mahasiswa)
// ====================================================================
func AddAttachment(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req models.Attachment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// get reference
	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	// ownership check
	if ref.StudentID.String() != userIDStr {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: not your achievement"})
	}

	// convert mongo id
	objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "invalid mongo id"})
	}

	// push attachment
	if err := repositories.AddAttachmentToAchievement(ctx, objID, req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to add attachment"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "attachment added"})
}

// ====================================================================
// GET /api/v1/achievements/:id — Detail (any allowed with proper role/owner)
// ====================================================================
func GetAchievementDetail(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "invalid mongo id"})
	}

	mongoDoc, err := repositories.GetAchievementMongoByID(ctx, objID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch mongo document"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"reference": ref,
			"detail":    mongoDoc,
		},
	})
}

// ====================================================================
// GET /api/v1/achievements — List (Mahasiswa: own, Dosen: advisees, Admin: all)
// ====================================================================
func ListAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	role, _ := c.Locals("role").(string)
	userIDStr, _ := c.Locals("userID").(string)

	if role == "Mahasiswa" {
		// get student profile to obtain student.id (students table)
		student, err := repositories.GetStudentByUserID(ctx, userIDStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch student profile"})
		}
		refs, err := repositories.GetAchievementReferencesByStudentIDs(ctx, []string{student.ID.String()})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch achievements"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": refs})
	}

	if role == "Dosen" {
		// get lecturer profile
		lecturer, err := repositories.GetLecturerByUserID(ctx, userIDStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch lecturer profile"})
		}
		// get students by advisor = lecturer.id
		students, err := repositories.GetStudentsByAdvisor(ctx, lecturer.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch advisees"})
		}
		// build studentIDs
		ids := []string{}
		for _, s := range students {
			ids = append(ids, s.ID.String())
		}
		refs, err := repositories.GetAchievementReferencesByStudentIDs(ctx, ids)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch achievements"})
		}
		return c.JSON(fiber.Map{"status": "success", "data": refs})
	}

	// Admin or others: list all
	refs, err := repositories.ListAllAchievementReferences(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch achievements"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": refs})
}

// ====================================================================
// GET /api/v1/achievements/:id/history — simple history (return reference timeline)
// ====================================================================
func GetAchievementHistory(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	// For now, return the timestamps available in the reference as history
	history := fiber.Map{
		"status":      ref.Status,
		"createdAt":   ref.CreatedAt,
		"submittedAt": ref.SubmittedAt,
		"verifiedAt":  ref.VerifiedAt,
		"rejection":   ref.RejectionNote,
	}

	return c.JSON(fiber.Map{"status": "success", "data": history})
}
