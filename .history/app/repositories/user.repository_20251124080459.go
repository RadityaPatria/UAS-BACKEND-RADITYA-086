package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// ----------------------------------------
// Create : Menambahkan user baru
// ----------------------------------------
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// ----------------------------------------
// FindByID : Mengambil user berdasarkan ID
// ----------------------------------------
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error
	return &user, err
}

// ----------------------------------------
// FindByUsername : Cari user berdasarkan username
// ----------------------------------------
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "username = ?", username).Error
	return &user, err
}

// ----------------------------------------
// Update : Mengupdate data user
// ----------------------------------------
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// ----------------------------------------
// Delete : Menghapus user (hard delete)
// ----------------------------------------
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, "id = ?", id).Error
}

// ----------------------------------------
// List : Mengambil semua user
// ----------------------------------------
func (r *UserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}
