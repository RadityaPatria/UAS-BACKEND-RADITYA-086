package services

import (
    "context"
    "errors"
    "time"

    "UAS-backend/app/models"
    "UAS-backend/app/repositories"

    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// ========================================================================
// FR-003 : CREATE ACHIEVEMENT (Mahasiswa)
// ========================================================================
func CreateAchievement(ctx context.Context, studentID string, payload models.AchievementMongo) (*models.AchievementReference, error) {

    // 1️⃣ Insert ke MongoDB (status selalu draft saat dibuat)
    payload.Status = models.StatusDraft
    mongoID, err := repositories.CreateAchievementMongo(ctx, &payload)
    if err != nil {
        return nil, err
    }

    // 2️⃣ Insert ke PostgreSQL (reference)
    ref := models.AchievementReference{
        ID:                 uuid.New(),
        StudentID:          uuid.MustParse(studentID),
        MongoAchievementID: mongoID.Hex(),
        Status:             models.StatusDraft,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := repositories.CreateAchievementReference(ctx, &ref); err != nil {
        return nil, err
    }

    // 3️⃣ Return sesuai FR-003
    return &ref, nil
}



// ========================================================================
// FR-004 : SUBMIT ACHIEVEMENT (Draft → Submitted)
// ========================================================================
func SubmitAchievement(ctx context.Context, refID, studentID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.StudentID.String() != studentID {
        return errors.New("forbidden: cannot submit other's achievement")
    }

    if ref.Status != models.StatusDraft {
        return errors.New("only DRAFT achievements can be submitted")
    }

    now := time.Now()
    ref.SubmittedAt = &now

    err = repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusSubmitted)
    if err != nil {
        return err
    }

    // FR-004: TODO → Create notification untuk dosen wali

    return nil
}



// ========================================================================
// FR-005 : DELETE ACHIEVEMENT (Draft only)
// ========================================================================
func DeleteAchievement(ctx context.Context, refID, studentID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.StudentID.String() != studentID {
        return errors.New("cannot delete another user's achievement")
    }

    if ref.Status != models.StatusDraft {
        return errors.New("only draft achievements can be deleted")
    }

    // PostgreSQL soft delete
    if err := repositories.SoftDeleteAchievementReference(ctx, refID); err != nil {
        return err
    }

    // MongoDB soft delete
    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    return repositories.SoftDeleteAchievementMongo(ctx, mongoID)
}



// ========================================================================
// FR-006 : LIST ACHIEVEMENTS based on role
// ========================================================================
func ListAchievements(ctx context.Context, role string, userID string) ([]models.AchievementReference, error) {

    switch role {

    case "Admin":
        return repositories.ListAllAchievementReferences(ctx)

    case "Dosen Wali":
        lecturerUUID, err := uuid.Parse(userID)
        if err != nil {
            return nil, errors.New("invalid UUID for lecturer")
        }

        advisees, err := repositories.GetAdviseesByLecturerID(ctx, lecturerUUID)
        if err != nil {
            return nil, err
        }

        return repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)

    case "Mahasiswa":
        return repositories.GetAchievementReferencesByStudentIDs(ctx, []string{userID})

    default:
        return nil, errors.New("invalid role")
    }
}



// ========================================================================
// FR-007 : VERIFY ACHIEVEMENT (Dosen Wali)
// ========================================================================
func VerifyAchievement(ctx context.Context, refID string, lecturerID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.Status != models.StatusSubmitted {
        return errors.New("only SUBMITTED achievements can be verified")
    }

    if err := repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified); err != nil {
        return err
    }

    return repositories.SetVerifiedBy(ctx, refID, lecturerID)
}



// FR-008 : REJECT ACHIEVEMENT (Dosen Wali)
func RejectAchievement(ctx context.Context, refID, lecturerID, note string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.Status != models.StatusSubmitted {
        return errors.New("only SUBMITTED achievements can be rejected")
    }

    return repositories.SetAchievementRejectionNote(ctx, refID, note)
}



// DETAIL ACHIEVEMENT (FR tambahan dari tabel)
func GetAchievementDetail(ctx context.Context, refID string) (*models.AchievementMongo, *models.AchievementReference, error) {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return nil, nil, errors.New("reference not found")
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    ach, err := repositories.GetAchievementMongoByID(ctx, mongoID)
    if err != nil {
        return nil, nil, errors.New("mongo achievement not found")
    }

    return ach, ref, nil
}



// ========================================================================
// HISTORY (Status Log)
// ========================================================================
func GetAchievementHistory(ctx context.Context, refID string) (*models.AchievementReference, error) {
    return repositories.GetAchievementReferenceByID(ctx, refID)
}



// ========================================================================
// ADD ATTACHMENT (Upload)
// ========================================================================
func AddAttachment(ctx context.Context, refID string, att models.Attachment) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement reference not found")
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    return repositories.AddAttachmentToAchievement(ctx, mongoID, att)
}
