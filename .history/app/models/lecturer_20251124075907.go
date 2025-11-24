package models

import (
	"time"

	"github.com/google/uuid"
)

type Lecturer struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	LecturerID  string    `gorm:"type:varchar(20);unique;not null"`
	Department  string    `gorm:"type:varchar(100)"`

	CreatedAt time.Time
}
