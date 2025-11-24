package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"UAS-backend/config"
	"UAS-backend/database"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	database.ConnectPostgres(cfg)

	// Try simple PostgreSQL ping
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("PostgreSQL SQL DB error:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("PostgreSQL Ping error:", err)
	} else {
		fmt.Println("PostgreSQL Ping OK ✔")
	}

	// Connect to MongoDB
	database.ConnectMongo(cfg)

	// Try MongoDB ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.MongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB Ping error:", err)
	} else {
		fmt.Println("MongoDB Ping OK ✔")
	}

	fmt.Println("SEMUA DATABASE BERHASIL TERHUBUNG! 🔥")
}
