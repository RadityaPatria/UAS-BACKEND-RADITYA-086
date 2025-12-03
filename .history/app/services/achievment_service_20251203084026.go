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

//
// ------------------------------------------------------------
// 1. CREATE ACHIEVEMENT (Draft) — Mahasiswa
// ------------------------------------------------------------
func CreateAchievement(ctx context.Context, studentID uuid.UUID, input *models.AchievementMongo) (*models.AchievementMongo, error) {

	// 1. Create in MongoDB
	input.StudentID = studentID.String()

	mongoID, err := repositories.CreateAchievementMongo(ctx, input)
	if err != nil {
		return nil, err
	}

	// 2. Create PostgreSQL reference
	ref := &models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             models.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err = repositories.CreateAchievementReference(ctx, ref)
	if err != nil {
		return nil, err
	}

	return input, nil
}

//
// ------------------------------------------------------------
// 2. UPDATE ACHIEVEMENT — Hanya Mahasiswa & status draft
// ------------------------------------------------------------
func UpdateAchievement(ctx context.Context, refID string, studentID uuid.UUID, updateData bson.M) error {

	// Check reference
	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	if ref.StudentID != studentID {
		return errors.New("tidak boleh mengupdate prestasi milik orang lain")
	}

	if ref.Status != models.StatusDraft {
		return errors.New("hanya draft yang bisa diupdate")
	}

	// Convert mongo id
	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	return repositories.UpdateAchievementMongo(ctx, oid, updateData)
}

//
// ------------------------------------------------------------
// 3. DELETE ACHIEVEMENT — Soft delete (draft only)
// ------------------------------------------------------------
func DeleteAchievement(ctx context.Context, refID string, studentID uuid.UUID) error {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	if ref.StudentID != studentID {
		return errors.New("tidak boleh menghapus prestasi milik orang lain")
	}

	if ref.Status != models.StatusDraft {
		return errors.New("hanya draft yang dapat dihapus")
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	// soft delete mongo
	err = repositories.SoftDeleteAchievementMongo(ctx, oid)
	if err != nil {
		return err
	}

	// soft delete postgres
	return repositories.SoftDeleteAchievementReference(ctx, refID)
}

//
// ------------------------------------------------------------
// 4. SUBMIT PRESTASI (Mahasiswa) — ubah status draft → submitted
// ------------------------------------------------------------
func SubmitAchievement(ctx context.Context, refID string, studentID uuid.UUID) error {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	if ref.StudentID != studentID {
		return errors.New("tidak boleh submit prestasi milik orang lain")
	}

	if ref.Status != models.StatusDraft {
		return errors.New("hanya draft yang dapat disubmit")
	}

	return repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusSubmitted)
}

//
// ------------------------------------------------------------
// 5. VERIFY (Dosen) — status submitted → verified
// ------------------------------------------------------------
func VerifyAchievement(ctx context.Context, refID string, lecturerID uuid.UUID) error {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	if ref.Status != models.StatusSubmitted {
		return errors.New("hanya prestasi submitted yang dapat diverifikasi")
	}

	// Set verified_by manually
	now := time.Now()

	err = repositories.UpdateAchievementReferenceStatus(ctx, refID, models.StatusVerified)
	if err != nil {
		return err
	}

	// UPDATE verified_by field
	_, err = repositories.DB.Exec(
		ctx,
		`UPDATE achievement_references 
		 SET verified_by=$2, verified_at=$3 
		 WHERE id=$1`,
		refID, lecturerID, now,
	)

	return err
}

//
// ------------------------------------------------------------
// 6. REJECT (Dosen) — status submitted → rejected + note
// ------------------------------------------------------------
func RejectAchievement(ctx context.Context, refID string, lecturerID uuid.UUID, note string) error {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	if ref.Status != models.StatusSubmitted {
		return errors.New("hanya submitted yang bisa ditolak")
	}

	return repositories.SetAchievementRejectionNote(ctx, refID, note)
}

//
// ------------------------------------------------------------
// 7. GET BY ROLE
// ------------------------------------------------------------

// Mahasiswa: hanya miliknya sendiri
func GetAchievementsForStudent(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error) {
	return repositories.GetAchievementReferencesByStudentIDs(ctx, []string{studentID.String()})
}

// Dosen wali: semua mahasiswa bimbingannya
func GetAchievementsForLecturer(ctx context.Context, lecturerID uuid.UUID) ([]models.AchievementReference, error) {

	students, err := repositories.GetStudentsByAdvisor(ctx, lecturerID.String())
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, s := range students {
		list = append(list, s.ID.String())
	}

	return repositories.GetAchievementReferencesByStudentIDs(ctx, list)
}

// Admin: semua
func GetAchievementsForAdmin(ctx context.Context) ([]models.AchievementReference, error) {
	return repositories.ListAllAchievementReferences(ctx)
}

//
// ------------------------------------------------------------
// 8. DETAIL
// ------------------------------------------------------------
func GetAchievementDetail(ctx context.Context, refID string) (*models.AchievementMongo, *models.AchievementReference, error) {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return nil, nil, err
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	mongoData, err := repositories.GetAchievementMongoByID(ctx, oid)
	if err != nil {
		return nil, nil, err
	}

	return mongoData, ref, nil
}

//
// ------------------------------------------------------------
// 9. UPLOAD ATTACHMENT
// ------------------------------------------------------------
func AddAchievementAttachment(ctx context.Context, refID string, att models.Attachment) error {

	ref, err := repositories.GetAchievementReferenceByID(ctx, refID)
	if err != nil {
		return err
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	return repositories.AddAttachmentToAchievement(ctx, oid, att)
}

//
// ------------------------------------------------------------
// 10. HISTORY
// (If needed, but we use reference timestamps)
// ------------------------------------------------------------
func GetAchievementHistory(ctx context.Context, refID string) (*models.AchievementReference, error) {
	return repositories.GetAchievementReferenceByID(ctx, refID)
}
