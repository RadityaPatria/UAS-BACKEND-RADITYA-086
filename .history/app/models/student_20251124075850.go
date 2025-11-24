package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	StudentID     string    `gorm:"type:varchar(20);unique;not null"`
	ProgramStudy  string    `gorm:"type:varchar(100)"`
	AcademicYear  string    `gorm:"type:varchar(10)"`
	AdvisorID     uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time
}
