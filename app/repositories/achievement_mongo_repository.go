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

// -------------------------------------------------------
// CreateAchievementMongo
// -------------------------------------------------------
func CreateAchievementMongo(ctx context.Context, a *models.AchievementMongo) (primitive.ObjectID, error) {
    a.ID = primitive.NewObjectID()
    a.CreatedAt = time.Now()
    a.UpdatedAt = time.Now()

    // WAJIB: default attachments
    if a.Attachments == nil {
        a.Attachments = []models.Attachment{}
    }

    _, err := collAchievements().InsertOne(ctx, a)
    return a.ID, err
}



// -------------------------------------------------------
// GetAchievementMongoByID
// -------------------------------------------------------
func GetAchievementMongoByID(ctx context.Context, id primitive.ObjectID) (*models.AchievementMongo, error) {
	var res models.AchievementMongo
	err := collAchievements().FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return &res, err
}

// -------------------------------------------------------
func UpdateAchievementMongo(ctx context.Context, id primitive.ObjectID, update bson.M) error {

    _, err := collAchievements().UpdateOne(
        ctx,
        bson.M{"_id": id},
        update,
    )

    return err
}


// -------------------------------------------------------
// SoftDeleteAchievementMongo (uses status=deleted)
// -------------------------------------------------------
func SoftDeleteAchievementMongo(ctx context.Context, id primitive.ObjectID) error {
	_, err := collAchievements().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":    models.StatusDeleted,
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

// -------------------------------------------------------
// AddAttachmentToAchievement
func AddAttachmentToAchievement(ctx context.Context, id primitive.ObjectID, att models.Attachment) error {
    att.UploadedAt = time.Now()

    // 1️⃣ Pastikan attachments sudah array
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

    // 2️⃣ Sekarang push aman
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



