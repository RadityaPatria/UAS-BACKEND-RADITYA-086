package main

import (
	"log"

	"UAS-backend/config"
	"UAS-backend/database"
	"UAS-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load .env
	cfg := config.LoadConfig()

	// Connect databases
	database.ConnectPostgres(cfg)
	database.ConnectMongo(cfg)

	// Create Fiber App
	app := fiber.New()

	// Register routes
	routes.SetupRoutes(app)

	// Start server
	log.Println("🚀 Server is running on port", cfg.AppPort)
	app.Listen(":" + cfg.AppPort)
}
