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

    // ambil body
    var update map[string]interface{}
    if err := c.BodyParser(&update); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // ambil reference di Postgres
    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // validasi mahasiswa pemilik
    studentLocal := c.Locals("studentID")
    if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    // hanya draft yang bisa diupdate
    if ref.Status != models.StatusDraft {
        return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
    }

    // update ke Mongo
    mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid mongo id"})
    }

    updateBson := bson.M{
        "$set": bson.M{
            "updatedAt": time.Now(),
        },
    }

    // merge user updates
    for k, v := range update {
        updateBson["$set"].(bson.M)[k] = v
    }

    if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateBson); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // update updated_at di Postgres TANPA ubah status
    if err := repositories.TouchAchievementReference(ctx, refID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

    studentLocal := c.Locals("studentID")
    if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    if ref.Status != models.StatusDraft {
        return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
    }

    // Update Postgres — FIX SUBMITTED_AT
    if err := repositories.SetSubmitted(ctx, refID); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "failed to update postgres: " + err.Error(),
        })
    }

    // Update Mongo
    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

    updateMongo := bson.M{
        "$set": bson.M{
            "status":      models.StatusSubmitted,
            "submittedAt": time.Now(),
            "updatedAt":   time.Now(),
        },
    }

    if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateMongo); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "failed to update mongo: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "submitted",
    })
}




// DELETE /api/v1/achievements/:id
func DeleteAchievement(c *fiber.Ctx) error {
    ctx := context.Background()
    refID := c.Params("id")

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "not found"})
    }

    role := c.Locals("role")

    // Mahasiswa boleh hapus jika status draft dan miliknya
    if role == "Mahasiswa" {
        studentLocal := c.Locals("studentID")
        if studentLocal == nil || ref.StudentID.String() != studentLocal.(string) {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }

        if ref.Status != models.StatusDraft {
            return c.Status(400).JSON(fiber.Map{"error": "only draft can be deleted"})
        }
    }

    // Admin bisa hapus apapun
    // Tidak perlu cek status dan kepemilikan

    // 1️⃣ Delete Postgres
    if err := repositories.SetDeleted(ctx, refID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // 2️⃣ Delete Mongo
    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err := repositories.SoftDeleteAchievementMongo(ctx, mongoID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "deleted",
    })
}




// GET /api/v1/achievements
func ListAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)

	var list []models.AchievementReference
	var err error

	switch role {
	case "Admin":
		list, err = repositories.ListAllAchievementReferences(ctx)

	case "Dosen Wali":
		lecturerID := c.Locals("lecturerID").(string)

        // list mahasiswa bimbingan
		advisees, err2 := repositories.GetAdviseesByLecturerID(ctx, uuid.MustParse(lecturerID))
		if err2 != nil {
			return c.Status(500).JSON(fiber.Map{"error": err2.Error()})
		}

		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)

	case "Mahasiswa":
		studentID := c.Locals("studentID")
		if studentID == nil {
			return c.Status(400).JSON(fiber.Map{"error": "student mapping not found"})
		}
		list, err = repositories.GetAchievementReferencesByStudentIDs(ctx, []string{studentID.(string)})

	default:
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": list})
}


func VerifyAchievement(c *fiber.Ctx) error {
    ctx := context.Background()
    refID := c.Params("id")

    lecturerIDraw := c.Locals("lecturerID")
    if lecturerIDraw == nil {
        return c.Status(403).JSON(fiber.Map{"error": "lecturerID not found"})
    }

    lecturerUUID, err := uuid.Parse(lecturerIDraw.(string))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid lecturer uuid"})
    }

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement reference not found"})
    }

    // 0️⃣ CEK apakah mahasiswa ini adalah anak bimbingan dosen tsb
    isAdvisee, err := repositories.IsStudentAdvisee(ctx, ref.StudentID.String(), lecturerUUID.String())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "advisor lookup failed"})
    }
    if !isAdvisee {
        return c.Status(403).JSON(fiber.Map{"error": "you are not the advisor for this student"})
    }

    if ref.Status != models.StatusSubmitted {
        return c.Status(400).JSON(fiber.Map{"error": "only submitted can be verified"})
    }

    if err := repositories.SetVerifiedBy(ctx, refID, lecturerUUID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to set verified_by: " + err.Error()})
    }

    if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to update status: " + err.Error()})
    }

    mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "invalid mongo ObjectID: " + err.Error()})
    }

    updateMongo := bson.M{
        "$set": bson.M{
            "status":     models.StatusVerified,
            "verifiedBy": lecturerUUID.String(),
            "verifiedAt": time.Now(),
            "updatedAt":  time.Now(),
        },
    }

    if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateMongo); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "mongo sync failed: " + err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "verified"})
}

// POST /api/v1/achievements/:id/reject  (Dosen Wali)
func RejectAchievement(c *fiber.Ctx) error {
    ctx := context.Background()
    refID := c.Params("id")

    // Ambil lecturerID dari JWT middleware
    lecturerIDraw := c.Locals("lecturerID")
    if lecturerIDraw == nil {
        return c.Status(403).JSON(fiber.Map{"error": "lecturerID not found"})
    }

    lecturerUUID, err := uuid.Parse(lecturerIDraw.(string))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid lecturer UUID"})
    }

    // Parsing body (harus ada note)
    var body struct {
        Note string `json:"note"`
    }
    if err := c.BodyParser(&body); err != nil || body.Note == "" {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body or empty note"})
    }

    // Ambil reference dari Postgres
    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement reference not found"})
    }

    // ❗ 1️⃣ Cek apakah mahasiswa ini adalah ANAK BIMBINGAN dosen
    isAdvisee, err := repositories.IsStudentAdvisee(ctx, ref.StudentID.String(), lecturerUUID.String())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "advisor lookup failed"})
    }
    if !isAdvisee {
        return c.Status(403).JSON(fiber.Map{
            "error": "you are not the advisor for this student",
        })
    }

    // Hanya status submitted yang boleh ditolak
    if ref.Status != models.StatusSubmitted {
        return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
    }

    // 2️⃣ UPDATE POSTGRES SESUAI STRUCT (status + rejection_note)
    if err := repositories.RejectAchievementPostgres(ctx, refID, body.Note); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "postgres update failed: " + err.Error(),
        })
    }

    // 3️⃣ UPDATE MONGO
    mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "invalid mongo ObjectID: " + err.Error(),
        })
    }

    updateMongo := bson.M{
        "$set": bson.M{
            "status":        models.StatusRejected,
            "rejectionNote": body.Note,
            "rejectedBy":    lecturerUUID.String(), 
            "rejectedAt":    time.Now(),
            "updatedAt":     time.Now(),
        },
    }

    if err := repositories.UpdateAchievementMongo(ctx, mongoID, updateMongo); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "mongo sync failed: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "rejected",
    })
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

func AddAttachment(c *fiber.Ctx) error {
    ctx := context.Background()
    refID := c.Params("id")

    // Ambil file dari form-data
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "file required"})
    }

    // Simpan file ke folder uploads
    filePath := "./uploads/" + file.Filename
    if err := c.SaveFile(file, filePath); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to save file"})
    }

    // Buat attachment sesuai model
    att := models.Attachment{
        FileName:   file.Filename,
        FileURL:    filePath,
        FileType:   file.Header.Get("Content-Type"),
        UploadedAt: time.Now(),
    }

    // Ambil reference dari Postgres
    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // Convert Mongo ID
    mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid mongo id"})
    }

    // Push attachment ke Mongo
    if err := repositories.AddAttachmentToAchievement(ctx, mongoID, att); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "attachment added",
        "attachment": att,
    })
}


