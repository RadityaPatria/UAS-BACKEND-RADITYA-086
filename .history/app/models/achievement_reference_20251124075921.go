package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusDraft     = "draft"
	StatusSubmitted = "submitted"
	StatusVerified  = "verified"
	StatusRejected  = "rejected"
)

type AchievementReference struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey"`
	StudentID          uuid.UUID  `gorm:"type:uuid;not null"`
	MongoAchievementID string     `gorm:"type:varchar(24);not null"`

	Status        string     `gorm:"type:varchar(20);not null;default:'draft'"`
	SubmittedAt   *time.Time
	VerifiedAt    *time.Time
	VerifiedBy    *uuid.UUID `gorm:"type:uuid"`
	RejectionNote *string

	CreatedAt time.Time
	UpdatedAt time.Time

	// ❗SOFT DELETE KHUSUS FR-005
	DeletedAt *time.Time `gorm:"index"`
}
