package main

import (
	"log"

	"UAS-backend/config"
	"UAS-backend/database"
	"UAS-backend/routes"

	_ "UAS-backend/docs" 

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title           UAS Backend API
// @version         1.0
// @description     API Sistem Prestasi Mahasiswa
// @termsOfService  http://swagger.io/terms/

// @contact.name   
// @contact.email  

// @host      localhost:3000
// @BasePath  /api/v1
// @schemes   http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	database.ConnectPostgres(cfg)
	database.ConnectMongo(cfg)

	// Create app
	app := fiber.New()

	// ✅ SWAGGER ROUTE (INI YANG TADI KURANG)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes
	routes.RegisterAuthRoutes(app)
	routes.RegisterUserRoutes(app)
	routes.RegisterAchievementRoutes(app)
	routes.RegisterStudentRoutes(app)
	routes.RegisterLecturerRoutes(app)
	routes.RegisterReportRoutes(app)

	log.Println("🚀 Server running on port", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
