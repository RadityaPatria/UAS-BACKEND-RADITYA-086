package database

import (
	"context"
	"log"
	"time"

	"UAS-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("❌ Failed to connect MongoDB:", err)
	}

	MongoClient = client
	MongoDB = client.Database(cfg.MongoDB)

	log.Println("✅ Connected to MongoDB successfully.")
}
