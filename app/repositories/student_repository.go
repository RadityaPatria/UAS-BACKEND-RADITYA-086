package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type StudentRepository struct {
	DB *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

// ----------------------------------------
// Create : Menambahkan student baru
// ----------------------------------------
func (r *StudentRepository) Create(s *models.Student) error {
	return r.DB.Create(s).Error
}

// ----------------------------------------
// FindByID : Ambil student berdasarkan ID
// ----------------------------------------
func (r *StudentRepository) FindByID(id uuid.UUID) (*models.Student, error) {
	var s models.Student
	err := r.DB.First(&s, "id = ?", id).Error
	return &s, err
}

// ----------------------------------------
// List : Ambil semua student
// ----------------------------------------
func (r *StudentRepository) List() ([]models.Student, error) {
	var s []models.Student
	err := r.DB.Find(&s).Error
	return s, err
}
