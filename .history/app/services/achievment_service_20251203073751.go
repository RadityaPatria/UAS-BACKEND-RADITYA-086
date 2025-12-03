package services

import (
	"context"
	"errors"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct{}

func NewAchievementService() *AchievementService {
	return &AchievementService{}
}

//
// ﹥ 1. CREATE ACHIEVEMENT (Mahasiswa)
//     Simpan data dinamis ke Mongo + buat reference di Postgre
//
func (s *AchievementService) CreateAchievement(ctx context.Context, studentID string, input models.AchievementMongo) (string, error) {
	// 1) Insert Mongo
	mongoID, err := repositories.CreateAchievementMongo(ctx, &input)
	if err != nil {
		return "", err
	}

	// 2) Build reference Postgre
	ref := models.AchievementReference{
		ID:                 models.GenerateULID(),
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             models.StatusDraft,
	}

	// 3) Insert Postgre
	if err := repositories.CreateAchievementReference(ctx, &ref); err != nil {
		return "", err
	}

	return ref.ID, nil
}

//
// ﹥ 2. GET DETAIL ACHIEVEMENT (Dapatkan dari both DB)
//
func (s *AchievementService) GetAchievementDetail(ctx context.Context, referenceID string) (*models.AchievementMongo, *models.AchievementReference, error) {

	// Get reference postgre
	ref, err := repositories.GetAchievementReferenceByID(ctx, referenceID)
	if err != nil {
		return nil, nil, err
	}

	// Convert mongoID string → ObjectID
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	// Get Mongo document
	ach, err := repositories.GetAchievementMongoByID(ctx, mongoID)
	if err != nil {
		return nil, nil, err
	}

	return ach, ref, nil
}

//
// ﹥ 3. SUBMIT ACHIEVEMENT (Mahasiswa)
//
func (s *AchievementService) SubmitAchievement(ctx context.Context, referenceID string) error {
	return repositories.UpdateAchievementReferenceStatus(ctx, referenceID, models.StatusSubmitted)
}

//
// ﹥ 4. VERIFY ACHIEVEMENT (Dosen/Admin)
//
func (s *AchievementService) VerifyAchievement(ctx context.Context, referenceID string, verifiedBy string) error {
	err := repositories.UpdateAchievementReferenceStatus(ctx, referenceID, models.StatusVerified)
	if err != nil {
		return err
	}

	// Set verified_by
	return repositories.SetVerifiedBy(ctx, referenceID, verifiedBy)
}

//
// ﹥ 5. REJECT ACHIEVEMENT (Dosen)
//
func (s *AchievementService) RejectAchievement(ctx context.Context, referenceID string, note string) error {
	if note == "" {
		return errors.New("rejection note required")
	}

	return repositories.SetAchievementRejectionNote(ctx, referenceID, note)
}

//
// ﹥ 6. ADD ATTACHMENT (Mongo)
//
func (s *AchievementService) AddAttachment(ctx context.Context, referenceID string, attachment models.Attachment) error {

	// find reference first
	ref, err := repositories.GetAchievementReferenceByID(ctx, referenceID)
	if err != nil {
		return err
	}

	// convert mongoID
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	return repositories.AddAttachmentToAchievement(ctx, mongoID, attachment)
}

//
// ﹥ 7. LIST ACHIEVEMENT BY Student IDs (Admin / Dosen)
//
func (s *AchievementService) ListAchievementsForStudents(ctx context.Context, studentIDs []string) ([]models.AchievementReference, error) {
	return repositories.GetAchievementReferencesByStudentIDs(ctx, studentIDs)
}
