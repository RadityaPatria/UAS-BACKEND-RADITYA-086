package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"type:varchar(50);unique;not null"`
	Email        string    `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	FullName     string    `gorm:"type:varchar(100);not null"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null"`
	IsActive     bool      `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
