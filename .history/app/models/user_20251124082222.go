package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	RoleID       uuid.UUID
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
