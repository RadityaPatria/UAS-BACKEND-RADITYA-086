package repositories

import (
	"UAS-backend/app/models"

	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	DB *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{DB: db}
}

// ----------------------------------------
// Assign : Memberi hak permission ke role
// ----------------------------------------
func (r *RolePermissionRepository) Assign(rp *models.RolePermission) error {
	return r.DB.Create(rp).Error
}

// ----------------------------------------
// GetPermissionsByRole : Ambil semua permission berdasarkan role
// ----------------------------------------
func (r *RolePermissionRepository) GetPermissionsByRole(roleID string) ([]models.RolePermission, error) {
	var list []models.RolePermission
	err := r.DB.Where("role_id = ?", roleID).Find(&list).Error
	return list, err
}
