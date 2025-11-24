package repositories

import (
	"UAS-backend/app/models"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type AchievementReferenceRepository struct {
	DB *gorm.DB
}

func NewAchievementReferenceRepository(db *gorm.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{DB: db}
}

// ----------------------------------------
// Create : Membuat reference prestasi baru
// ----------------------------------------
func (r *AchievementReferenceRepository) Create(ar *models.AchievementReference) error {
	return r.DB.Create(ar).Error
}

// ----------------------------------------
// Update : Mengupdate reference prestasi
// ----------------------------------------
func (r *AchievementReferenceRepository) Update(ar *models.AchievementReference) error {
	return r.DB.Save(ar).Error
}

// ----------------------------------------
// SoftDelete : Menghapus prestasi (tanpa menghilangkan data)
// ----------------------------------------
func (r *AchievementReferenceRepository) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.DB.Model(&models.AchievementReference{}).
		Where("id = ?", id).
		Update("deleted_at", &now).Error
}

// ----------------------------------------
// FindByID : Ambil reference prestasi berdasarkan ID
// ----------------------------------------
func (r *AchievementReferenceRepository) FindByID(id uuid.UUID) (*models.AchievementReference, error) {
	var ar models.AchievementReference
	err := r.DB.First(&ar, "id = ?", id).Error
	return &ar, err
}
