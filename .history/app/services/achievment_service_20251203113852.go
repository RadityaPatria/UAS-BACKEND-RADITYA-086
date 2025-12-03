package services

import (
    "context"
    "errors"
    "time"

    "UAS-backend/app/models"
    "UAS-backend/app/repositories"

    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// ======================================================
// 1. CREATE ACHIEVEMENT  (Mahasiswa)
// ======================================================
func CreateAchievement(ctx context.Context, studentID string, payload models.AchievementMongo) (*models.AchievementReference, error) {

    // Insert Mongo first
    mongoID, err := repositories.CreateAchievementMongo(ctx, &payload)
    if err != nil {
        return nil, err
    }

    // Insert reference (Postgres)
    ref := models.AchievementReference{
        ID:                 uuid.New(),
        StudentID:          uuid.MustParse(studentID),
        MongoAchievementID: mongoID.Hex(),
        Status:             models.StatusDraft,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    err = repositories.CreateAchievementReference(ctx, &ref)
    if err != nil {
        return nil, err
    }

    return &ref, nil
}


// ======================================================
// 2. UPDATE ACHIEVEMENT (Hanya Draft, Mahasiswa)
// ======================================================
func UpdateAchievement(ctx context.Context, refID, studentID string, update bson.M) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.StudentID.String() != studentID {
        return errors.New("forbidden: cannot edit others' achievement")
    }

    if ref.Status != models.StatusDraft {
        return errors.New("only draft achievements can be updated")
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

    return repositories.UpdateAchievementMongo(ctx, mongoID, update)
}


// ======================================================
// 3. SUBMIT ACHIEVEMENT (Mahasiswa)
// ======================================================
func SubmitAchievement(ctx context.Context, refID, studentID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.StudentID.String() != studentID {
        return errors.New("forbidden: cannot submit others’ achievement")
    }

    if ref.Status != models.StatusDraft {
        return errors.New("only draft achievements can be submitted")
    }

    now := time.Now()
    ref.SubmittedAt = &now

    err = repositories.UpdateAchievementReferenceStatus(ctx, ref.ID.String(), models.StatusSubmitted)
    return err
}


// ======================================================
// 4. VERIFY (Approved by Dosen Wali)
// ======================================================
func VerifyAchievement(ctx context.Context, refID string, lecturerID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.Status != models.StatusSubmitted {
        return errors.New("only submitted achievements can be verified")
    }

    // Set status verified
    err = repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified)
    if err != nil {
        return err
    }

    // Set who verified
    return repositories.SetVerifiedBy(ctx, refID, lecturerID)
}


// ======================================================
// 5. REJECT (Dosen Wali)
// ======================================================
func RejectAchievement(ctx context.Context, refID, lecturerID, note string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.Status != models.StatusSubmitted {
        return errors.New("only submitted achievements can be rejected")
    }

    // Save rejection note + change status
    return repositories.SetAchievementRejectionNote(ctx, refID, note)
}


// ======================================================
// 6. DELETE (Soft Delete - Mahasiswa only on draft)
// ======================================================
func DeleteAchievement(ctx context.Context, refID, studentID string) error {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return errors.New("achievement not found")
    }

    if ref.StudentID.String() != studentID {
        return errors.New("forbidden: cannot delete others’ achievement")
    }

    if ref.Status != models.StatusDraft {
        return errors.New("only draft achievements can be deleted")
    }

    // Soft delete in Postgres
    err = repositories.SoftDeleteAchievementReference(ctx, refID)
    if err != nil {
        return err
    }

    // Soft delete in Mongo
    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    return repositories.SoftDeleteAchievementMongo(ctx, mongoID)
}


// ======================================================
// 7. DETAIL ACHIEVEMENT
// ======================================================
func GetAchievementDetail(ctx context.Context, refID string) (*models.AchievementMongo, *models.AchievementReference, error) {

    ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
    if err != nil {
        return nil, nil, errors.New("achievement reference not found")
    }

    mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    ach, err := repositories.GetAchievementMongoByID(ctx, mongoID)
    if err != nil {
        return nil, nil, errors.New("achievement mongo not found")
    }

    return ach, ref, nil
}


// ======================================================
// 8. LIST ACHIEVEMENTS BASED ON ROLE
// ======================================================
func ListAchievements(ctx context.Context, role string, userID string) ([]models.AchievementReference, error) {

    switch role {

    case "Admin":
        return repositories.ListAllAchievementReferences(ctx)

    case "Dosen Wali":

        // Convert string → UUID
        lecturerUUID, err := uuid.Parse(userID)
        if err != nil {
            return nil, errors.New("invalid lecturer UUID in token")
        }

        // Get all advisees (students under this lecturer)
        advisees, err := repositories.GetAdviseesByLecturerID(ctx, lecturerUUID)
        if err != nil {
            return nil, err
        }

        // Get all achievement references of those students
        return repositories.GetAchievementReferencesByStudentIDs(ctx, advisees)

    case "Mahasiswa":
        return repositories.GetAchievementReferencesByStudentIDs(ctx, []string{userID})

    default:
        return nil, errors.New("invalid role")
    }
}

