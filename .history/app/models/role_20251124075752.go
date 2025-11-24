package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(50);unique;not null"`
	Description string    `gorm:"type:text"`

	CreatedAt time.Time
}
