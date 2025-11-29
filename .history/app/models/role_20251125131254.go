package models

import "time"
import "github.com/google/uuid"

type Role struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
}
