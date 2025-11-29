package main

import (
	"log"

	"UAS-backend/config"
	"UAS-backend/database"
	"UAS-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect Databases
	database.ConnectPostgres(cfg)
	database.ConnectMongo(cfg)

	// Create Fiber App
	app := fiber.New()

	// Register Routes
	routes.RegisterAuthRoutes(app)
	routes.RegisterUserRoutes(app) // <-- PENTING! Daftarkan User Routes

	// Start server
	log.Println("🚀 Server running on port", cfg.AppPort)
	app.Listen(":" + cfg.AppPort)
}
