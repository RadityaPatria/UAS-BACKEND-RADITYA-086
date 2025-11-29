package main

import (
	"log"

	"UAS-backend/config"
	"UAS-backend/database"
	"UAS-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load .env / config
	cfg := config.LoadConfig()

	// Connect Databases
	database.ConnectPostgres(cfg)
	database.ConnectMongo(cfg)

	// Init Fiber
	app := fiber.New()

	// Register all routes
	api := app.Group("/api/v1")     // <– PREFIX GLOBAL
	routes.RegisterAuthRoutes(api)  // <– seluruh auth ditempel ke /api/v1/auth

	log.Println("🚀 Server running on port", cfg.AppPort)
	app.Listen(":" + cfg.AppPort)
}
