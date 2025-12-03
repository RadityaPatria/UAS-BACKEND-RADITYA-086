package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

type AchievementMongo struct {
    ID             primitive.ObjectID       `bson:"_id" json:"id"`
    AchievementType string                  `bson:"achievementType" json:"achievementType"`
    Title          string                   `bson:"title" json:"title"`
    Description    string                   `bson:"description" json:"description"`
    Details        map[string]interface{}   `bson:"details" json:"details"`
    Tags           []string                 `bson:"tags" json:"tags"`
    Status         string                   `bson:"status" json:"status"`
    StudentID      string                   `bson:"studentId" json:"studentId"`
    CreatedAt      time.Time                `bson:"createdAt" json:"createdAt"`
    UpdatedAt      time.Time                `bson:"updatedAt" json:"updatedAt"`
}

