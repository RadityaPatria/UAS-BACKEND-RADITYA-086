package main

import (
	"log"

	"UAS-backend/config"
	"UAS-backend/database"
	"UAS-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	database.ConnectPostgres(cfg)
	database.ConnectMongo(cfg)

	app := fiber.New()
	routes.RegisterAuthRoutes(app)

	log.Println("🚀 Server running on port", cfg.AppPort)
	app.Listen(":" + cfg.AppPort)
}
