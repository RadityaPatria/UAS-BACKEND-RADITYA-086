package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type LecturerRepository struct {
	DB *gorm.DB
}

func NewLecturerRepository(db *gorm.DB) *LecturerRepository {
	return &LecturerRepository{DB: db}
}

// ----------------------------------------
// Create : Menambahkan dosen baru
// ----------------------------------------
func (r *LecturerRepository) Create(l *models.Lecturer) error {
	return r.DB.Create(l).Error
}

// ----------------------------------------
// FindByID : Ambil dosen berdasarkan ID
// ----------------------------------------
func (r *LecturerRepository) FindByID(id uuid.UUID) (*models.Lecturer, error) {
	var l models.Lecturer
	err := r.DB.First(&l, "id = ?", id).Error
	return &l, err
}

// ----------------------------------------
// List : Ambil semua dosen
// ----------------------------------------
func (r *LecturerRepository) List() ([]models.Lecturer, error) {
	var l []models.Lecturer
	err := r.DB.Find(&l).Error
	return l, err
}
