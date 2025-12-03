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

// collAchievements -> helper collection handle
func collAchievements() *mongo.Collection {
	return database.MongoDB.Collection("achievements")
}

// CreateAchievementMongo -> simpan achievement dinamis di Mongo
func CreateAchievementMongo(ctx context.Context, a *models.AchievementMongo) (primitive.ObjectID, error) {
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := collAchievements().InsertOne(ctx, a)
	return a.ID, err
}

// GetAchievementMongoByID -> ambil achievement dari Mongo by ObjectID
func GetAchievementMongoByID(ctx context.Context, id primitive.ObjectID) (*models.AchievementMongo, error) {
	var res models.AchievementMongo
	err := collAchievements().FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return &res, err
}

// UpdateAchievementMongo -> update field di Mongo (partial update)
func UpdateAchievementMongo(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := collAchievements().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

// SoftDeleteAchievementMongo -> tandai deleted=true (soft delete)
func SoftDeleteAchievementMongo(ctx context.Context, id primitive.ObjectID) error {
	_, err := collAchievements().UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"deleted": true, "updatedAt": time.Now()}},
	)
	return err
}

// AddAttachmentToAchievement -> push attachment ke array attachments
func AddAttachmentToAchievement(ctx context.Context, id primitive.ObjectID, att models.Attachment) error {
	att.UploadedAt = time.Now()

	_, err := collAchievements().UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{
			"$push": bson.M{"attachments": att},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)

	return err
}
