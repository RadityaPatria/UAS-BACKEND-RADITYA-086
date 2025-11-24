package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

// ----------------------------------------
// Create : Menambahkan role baru
// ----------------------------------------
func (r *RoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}

// ----------------------------------------
// FindByID : Cari role berdasarkan ID
// ----------------------------------------
func (r *RoleRepository) FindByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.DB.First(&role, "id = ?", id).Error
	return &role, err
}

// ----------------------------------------
// FindByName : Cari role berdasarkan nama
// ----------------------------------------
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.DB.First(&role, "name = ?", name).Error
	return &role, err
}

// ----------------------------------------
// List : Mengambil semua role
// ----------------------------------------
func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.Find(&roles).Error
	return roles, err
}
