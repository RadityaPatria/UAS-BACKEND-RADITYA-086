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
	// Load config
	cfg := config.LoadConfig()

	// Connect PostgreSQL
	database.ConnectPostgres(cfg)

	// PostgreSQL Ping Test
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := database.DB.Ping(ctx); err != nil {
		log.Fatal("PostgreSQL Ping error:", err)
	} else {
		fmt.Println("PostgreSQL Ping OK ✔")
	}

	// Connect MongoDB
	database.ConnectMongo(cfg)

	// MongoDB Ping Test
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := database.MongoClient.Ping(ctx2, nil); err != nil {
		log.Fatal("MongoDB Ping error:", err)
	} else {
		fmt.Println("MongoDB Ping OK ✔")
	}

	fmt.Println("SEMUA DATABASE BERHASIL TERHUBUNG! 🔥")
}
