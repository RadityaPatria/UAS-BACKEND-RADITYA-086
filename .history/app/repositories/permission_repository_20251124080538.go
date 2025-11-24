package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

// ----------------------------------------
// Create : Menambahkan permission baru
// ----------------------------------------
func (r *PermissionRepository) Create(p *models.Permission) error {
	return r.DB.Create(p).Error
}

// ----------------------------------------
// FindByID : Cari permission berdasarkan ID
// ----------------------------------------
func (r *PermissionRepository) FindByID(id uuid.UUID) (*models.Permission, error) {
	var perm models.Permission
	err := r.DB.First(&perm, "id = ?", id).Error
	return &perm, err
}

// ----------------------------------------
// List : Mengambil semua permission
// ----------------------------------------
func (r *PermissionRepository) List() ([]models.Permission, error) {
	var perms []models.Permission
	err := r.DB.Find(&perms).Error
	return perms, err
}
