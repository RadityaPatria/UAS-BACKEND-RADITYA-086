package repository

import (
	"context"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAchievement(ctx context.Context, a *models.AchievementMongo) (primitive.ObjectID, error) {
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := database.MongoDB.Collection("achievements").InsertOne(ctx, a)
	return a.ID, err
}

func GetAchievementByID(ctx context.Context, id primitive.ObjectID) (*models.AchievementMongo, error) {
	var result models.AchievementMongo
	err := database.MongoDB.Collection("achievements").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&result)

	return &result, err
}

func UpdateAchievement(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()

	_, err := database.MongoDB.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})

	return err
}

func SoftDeleteAchievement(ctx context.Context, id primitive.ObjectID) error {
	_, err := database.MongoDB.Collection("achievements").
		UpdateOne(ctx, bson.M{"_id": id},
			bson.M{"$set": bson.M{"deleted": true, "updatedAt": time.Now()}})
	return err
}

func AddAttachment(ctx context.Context, id primitive.ObjectID, att models.Attachment) error {
	att.UploadedAt = time.Now()

	_, err := database.MongoDB.Collection("achievements").
		UpdateOne(ctx,
			bson.M{"_id": id},
			bson.M{
				"$push": bson.M{"attachments": att},
				"$set":  bson.M{"updatedAt": time.Now()},
			})
	return err
}
