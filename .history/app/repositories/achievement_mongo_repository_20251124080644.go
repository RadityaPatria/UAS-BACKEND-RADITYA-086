package repositories

import (
	"UAS-backend/app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongoRepository struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) *AchievementMongoRepository {
	return &AchievementMongoRepository{
		Collection: db.Collection("achievements"),
	}
}

// ----------------------------------------
// Create : Menambahkan prestasi ke MongoDB
// ----------------------------------------
func (r *AchievementMongoRepository) Create(ctx context.Context, data *models.AchievementMongo) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, data)
	return err
}

// ----------------------------------------
// FindByID : Mengambil prestasi MongoDB berdasarkan ID
// ----------------------------------------
func (r *AchievementMongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.AchievementMongo, error) {
	var result models.AchievementMongo
	err := r.Collection.FindOne(ctx, bson.M{"_id": id, "isDeleted": false}).Decode(&result)
	return &result, err
}

// ----------------------------------------
// SoftDelete : Menghapus prestasi di Mongo dengan flag
// ----------------------------------------
func (r *AchievementMongoRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"deletedAt": now,
		},
	}
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
