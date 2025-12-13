package repositories

import (
	"context"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func collAchievements() *mongo.Collection {
	return database.MongoDB.Collection("achievements")
}

// CreateAchievementMongo -> simpan detail prestasi ke MongoDB | FR-007
func CreateAchievementMongo(ctx context.Context, a *models.AchievementMongo) (primitive.ObjectID, error) {
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	// default attachments wajib ada
	if a.Attachments == nil {
		a.Attachments = []models.Attachment{}
	}

	_, err := collAchievements().InsertOne(ctx, a)
	return a.ID, err
}

// GetAchievementMongoByID -> ambil detail prestasi MongoDB | FR-007
func GetAchievementMongoByID(ctx context.Context, id primitive.ObjectID) (*models.AchievementMongo, error) {
	var res models.AchievementMongo
	err := collAchievements().FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return &res, err
}

// UpdateAchievementMongo -> update data prestasi MongoDB | FR-007
func UpdateAchievementMongo(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := collAchievements().UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)
	return err
}

// SoftDeleteAchievementMongo -> soft delete prestasi MongoDB | FR-009
func SoftDeleteAchievementMongo(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"status":    "deleted",
			"deletedAt": time.Now(),
		},
	}

	_, err := collAchievements().UpdateByID(ctx, id, update)
	return err
}

// AddAttachmentToAchievement -> tambah lampiran prestasi | FR-007
func AddAttachmentToAchievement(ctx context.Context, id primitive.ObjectID, att models.Attachment) error {
	att.UploadedAt = time.Now()

	// pastikan field attachments berupa array
	_, _ = collAchievements().UpdateOne(
		ctx,
		bson.M{"_id": id, "attachments": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"attachments": []models.Attachment{}}},
	)

	_, _ = collAchievements().UpdateOne(
		ctx,
		bson.M{"_id": id, "attachments": nil},
		bson.M{"$set": bson.M{"attachments": []models.Attachment{}}},
	)

	// push attachment
	_, err := collAchievements().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$push": bson.M{"attachments": att},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)

	return err
}
