package services

import (
	"context"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



// @Summary      Create Achievement
// @Description  Mahasiswa membuat prestasi (MongoDB + PostgreSQL)
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /achievements [post]
//
// CreateAchievement -> mahasiswa membuat prestasi (Mongo + Postgres) | FR-010
func CreateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()

	// Ambil studentID dari JWT middleware
	studentIDraw := c.Locals("studentID")
	userIDraw := c.Locals("userID")

	var studentID string
	if studentIDraw != nil {
		studentID = studentIDraw.(string)
	} else if userIDraw != nil {
		userID := userIDraw.(string)
		sid, err := repositories.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "student record not found for this user",
			})
		}
		studentID = sid
	} else {
		return c.Status(401).JSON(fiber.Map{"error": "unauthenticated"})
	}

	var payload models.AchievementMongo
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	payload.Status = models.StatusDraft
	payload.StudentID = studentID

	// Simpan ke MongoDB
	mongoID, err := repositories.CreateAchievementMongo(ctx, &payload)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Simpan reference ke PostgreSQL
	ref := models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.MustParse(studentID),
		MongoAchievementID: mongoID.Hex(),
		Status:             models.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := repositories.CreateAchievementReference(ctx, &ref); err != nil {
		_ = repositories.SoftDeleteAchievementMongo(ctx, mongoID)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"reference": ref,
			"mongo_id":  mongoID.Hex(),
		},
	})
}


// @Summary      Update Achievement
// @Description  Update prestasi (hanya status DRAFT)
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /achievements/{id} [put]
//
// UpdateAchievement -> update prestasi (hanya DRAFT) | FR-010
func UpdateAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	studentLocal := c.Locals("studentID")
	if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if ref.Status != models.StatusDraft {
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
	}

	mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid mongo id"})
	}

	updateBson := bson.M{
		"$set": bson.M{"updatedAt": time.Now()},
	}

	for k, v := range update {
		updateBson["$set"].(bson.M)[k] = v
	}

	if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateBson); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	_ = repositories.TouchAchievementReference(ctx, refID)

	return c.JSON(fiber.Map{"status": "success", "message": "updated"})
}

// @Summary      Submit Achievement
// @Description  Kirim prestasi ke dosen wali
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /achievements/{id}/submit [post]
//
// SubmitAchievement -> kirim prestasi ke dosen wali | FR-011
func SubmitAchievement(c *fiber.Ctx) error {
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
		return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
	}

	if err := repositories.SetSubmitted(ctx, refID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	_ = repositories.UpdateAchievementMongo(ctx, mongoID, bson.M{
		"$set": bson.M{
			"status":      models.StatusSubmitted,
			"submittedAt": time.Now(),
			"updatedAt":   time.Now(),
		},
	})

	return c.JSON(fiber.Map{"status": "success", "message": "submitted"})
}

// @Summary      Delete Achievement
// @Description  Hapus prestasi (Mahasiswa / Admin)
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /achievements/{id} [delete]
//
// DeleteAchievement -> hapus prestasi (Mahasiswa/Admin) | FR-012
func DeleteAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	role := c.Locals("role")

	if role == "Mahasiswa" {
		studentLocal := c.Locals("studentID")
		if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		if ref.Status != models.StatusDraft {
			return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
		}
	}

	_ = repositories.SetDeleted(ctx, refID)
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	_ = repositories.SoftDeleteAchievementMongo(ctx, mongoID)

	return c.JSON(fiber.Map{"status": "success", "message": "deleted"})
}


// @Summary      List Achievements
// @Description  List prestasi berdasarkan role (Admin / Dosen Wali / Mahasiswa)
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      403 {object} map[string]string
// @Router       /achievements [get]
//
// ListAchievements -> list prestasi berdasarkan role | FR-013
func ListAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)

	var list []models.AchievementReference
	var err error

	switch role {
	case "Admin":
		list, err = repositories.ListAllAchievementReferences(ctx)

	case "Dosen Wali":
		lecturerID := uuid.MustParse(c.Locals("lecturerID").(string))
		advisees, _ := repositories.GetAdviseesByLecturerID(ctx, lecturerID)
		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)

	case "Mahasiswa":
		studentID := c.Locals("studentID").(string)
		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, []string{studentID})

	default:
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": list})
}


// @Summary      Verify Achievement
// @Description  Verifikasi prestasi oleh Dosen Wali
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /achievements/{id}/verify [post]
//
// VerifyAchievement -> verifikasi prestasi oleh dosen wali | FR-014
func VerifyAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	lecturerUUID := uuid.MustParse(c.Locals("lecturerID").(string))
	ref, _ := repositories.GetAchievementReferenceByID(ctx, refID)

	if ok, _ := repositories.IsStudentAdvisee(ctx, ref.StudentID.String(), lecturerUUID.String()); !ok {
		return c.Status(403).JSON(fiber.Map{"error": "not advisor"})
	}

	_ = repositories.SetVerifiedBy(ctx, refID, lecturerUUID)
	_ = repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified)

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	_ = repositories.UpdateAchievementMongo(ctx, mongoID, bson.M{
		"$set": bson.M{
			"status":     models.StatusVerified,
			"verifiedBy": lecturerUUID.String(),
			"verifiedAt": time.Now(),
			"updatedAt":  time.Now(),
		},
	})

	return c.JSON(fiber.Map{"status": "success", "message": "verified"})
}


// @Summary      Reject Achievement
// @Description  Tolak prestasi oleh Dosen Wali
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /achievements/{id}/reject [post]
//
// RejectAchievement -> tolak prestasi oleh dosen wali | FR-014
func RejectAchievement(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)

	lecturerUUID := uuid.MustParse(c.Locals("lecturerID").(string))
	ref, _ := repositories.GetAchievementReferenceByID(ctx, refID)

	_ = repositories.RejectAchievementPostgres(ctx, refID, body.Note)

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	_ = repositories.UpdateAchievementMongo(ctx, mongoID, bson.M{
		"$set": bson.M{
			"status":        models.StatusRejected,
			"rejectionNote": body.Note,
			"rejectedBy":    lecturerUUID.String(),
			"rejectedAt":    time.Now(),
			"updatedAt":     time.Now(),
		},
	})

	return c.JSON(fiber.Map{"status": "success", "message": "rejected"})
}


// @Summary      Achievement Detail
// @Description  Detail prestasi (MongoDB + PostgreSQL)
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]string
// @Router       /achievements/{id} [get]
//
// GetAchievementDetail -> detail prestasi (Mongo + Postgres) | FR-013
func GetAchievementDetail(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, _ := repositories.GetAchievementReferenceByID(ctx, refID)
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	ach, _ := repositories.GetAchievementMongoByID(ctx, mongoID)

	return c.JSON(fiber.Map{
		"status":      "success",
		"reference":   ref,
		"achievement": ach,
	})
}


// @Summary      Achievement History
// @Description  Riwayat status prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Achievement Reference ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]string
// @Router       /achievements/{id}/history [get]
//
// GetAchievementHistory -> riwayat status prestasi | FR-013
func GetAchievementHistory(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	ref, _ := repositories.GetAchievementReferenceByID(ctx, refID)
	return c.JSON(fiber.Map{"status": "success", "history": ref})
}


// @Summary      Add Achievement Attachment
// @Description  Upload lampiran prestasi
// @Tags         Achievements
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path string true "Achievement Reference ID"
// @Param        file formData file true "Attachment File"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /achievements/{id}/attachment [post]
//
// AddAttachment -> upload lampiran prestasi | FR-010
func AddAttachment(c *fiber.Ctx) error {
	ctx := context.Background()
	refID := c.Params("id")

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}

	path := "./uploads/" + file.Filename
	_ = c.SaveFile(file, path)

	att := models.Attachment{
		FileName:   file.Filename,
		FileURL:    path,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	ref, _ := repositories.GetAchievementReferenceByID(ctx, refID)
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	_ = repositories.AddAttachmentToAchievement(ctx, mongoID, att)

	return c.JSON(fiber.Map{
		"status":     "success",
		"attachment": att,
	})
}
