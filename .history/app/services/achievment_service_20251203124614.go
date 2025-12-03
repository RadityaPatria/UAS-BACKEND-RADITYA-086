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


// ========================================================================
// FR-003 : CREATE ACHIEVEMENT (Mahasiswa)
// ========================================================================
func CreateAchievement(c *fiber.Ctx) error {

    studentID := c.Locals("userID").(string)

    var payload models.AchievementMongo
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
    }

    ctx := context.Background()

    payload.Status = models.StatusDraft

    mongoID, err := repositories.CreateAchievementMongo(ctx, &payload)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    ref := models.AchievementReference{
        ID:                 uuid.New(),
        StudentID:          uuid.MustParse(studentID),
        MongoAchievementID: mongoID.Hex(),
        Status:             models.StatusDraft,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := repositories.CreateAchievementReference(ctx, &ref); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   ref,
    })
}



// ========================================================================
// UPDATE ACHIEVEMENT (Mahasiswa - hanya draft)
// ========================================================================
func UpdateAchievement(c *fiber.Ctx) error {

    studentID := c.Locals("userID").(string)
    refID := c.Params("id")
    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    if ref.StudentID.String() != studentID {
        return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
    }

    if ref.Status != models.StatusDraft {
        return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
    }

    var payload models.AchievementMongo
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // convert struct → bson.M
    update := bson.M{
        "achievementType": payload.AchievementType,
        "title":           payload.Title,
        "description":     payload.Description,
        "details":         payload.Details,
        "tags":            payload.Tags,
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

    if err := repositories.UpdateAchievementMongo(ctx, mongoID, update); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    repositories.UpdateAchievementReferenceTime(ctx, refID)

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "achievement updated",
    })
}



// ========================================================================
// FR-004 : SUBMIT ACHIEVEMENT
// ========================================================================
func SubmitAchievement(c *fiber.Ctx) error {

    studentID := c.Locals("userID").(string)
    refID := c.Params("id")
    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    if ref.StudentID.String() != studentID {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    if ref.Status != models.StatusDraft {
        return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
    }

    now := time.Now()
    ref.SubmittedAt = &now

    if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusSubmitted); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "achievement submitted",
    })
}



// ========================================================================
// FR-005 : DELETE (Soft delete)
// ========================================================================
func DeleteAchievement(c *fiber.Ctx) error {

    studentID := c.Locals("userID").(string)
    refID := c.Params("id")

    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    if ref.StudentID.String() != studentID {
        return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
    }

    if ref.Status != models.StatusDraft {
        return c.Status(400).JSON(fiber.Map{"error": "cannot delete non-draft"})
    }

    repositories.SoftDeleteAchievementReference(ctx, refID)

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    repositories.SoftDeleteAchievementMongo(ctx, mongoID)

    return c.JSON(fiber.Map{"status": "success", "message": "deleted"})
}



// ========================================================================
// FR-006 : LIST ACHIEVEMENTS BY ROLE
// ========================================================================
func ListAchievements(c *fiber.Ctx) error {

    ctx := context.Background()
    role := c.Locals("role").(string)
    userID := c.Locals("userID").(string)

    var result []models.AchievementReference
    var err error

    switch role {

    case "Admin":
        result, err = repositories.ListAllAchievementReferences(ctx)

    case "Dosen Wali":
        lecturerUUID, err2 := uuid.Parse(userID)
        if err2 != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid lecturer UUID"})
        }

        advisees, err2 := repositories.GetAdviseesByLecturerID(ctx, lecturerUUID)
        if err2 != nil {
            return c.Status(500).JSON(fiber.Map{"error": err2.Error()})
        }

        result, err = repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)

    case "Mahasiswa":
        result, err = repositories.GetAchievementReferencesByStudentIDs(ctx, []string{userID})

    default:
        return c.Status(400).JSON(fiber.Map{"error": "invalid role"})
    }

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "data": result})
}



// ========================================================================
// FR-007 : VERIFY ACHIEVEMENT (Dosen Wali)
// ========================================================================
func VerifyAchievement(c *fiber.Ctx) error {

    lecturerID := c.Locals("userID").(string)
    refID := c.Params("id")

    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
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



// ========================================================================
// FR-008 : REJECT (Dosen Wali)
// ========================================================================
func RejectAchievement(c *fiber.Ctx) error {

    refID := c.Params("id")
    ctx := context.Background()

    var body struct {
        Note string `json:"note"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    if ref.Status != models.StatusSubmitted {
        return c.Status(400).JSON(fiber.Map{"error": "only submitted can be rejected"})
    }

    if err := repositories.SetAchievementRejectionNote(ctx, refID, body.Note); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "achievement rejected"})
}



// ========================================================================
// DETAIL ACHIEVEMENT
// ========================================================================
func GetAchievementDetail(c *fiber.Ctx) error {

    refID := c.Params("id")
    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    ach, err := repositories.GetAchievementMongoByID(ctx, mongoID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "mongo data not found"})
    }

    return c.JSON(fiber.Map{
        "status":      "success",
        "reference":   ref,
        "achievement": ach,
    })
}



// ========================================================================
// HISTORY (Status Log)
// ========================================================================
func GetAchievementHistory(c *fiber.Ctx) error {

    refID := c.Params("id")
    ctx := context.Background()

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "not found"})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "history": ref,
    })
}



// ========================================================================
// ADD ATTACHMENT
// ========================================================================
func AddAttachment(c *fiber.Ctx) error {

    refID := c.Params("id")
    ctx := context.Background()

    var att models.Attachment
    if err := c.BodyParser(&att); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

    if err := repositories.AddAttachmentToAchievement(ctx, mongoID, att); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "attachment added",
    })
}
