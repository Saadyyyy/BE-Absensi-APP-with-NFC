package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"myapp/config"
	"myapp/models"
	"myapp/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Auto-migrate database tables
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.School{},
		&models.Attendance{},
		&models.Student{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Setup routes
	routes.SetupRoutes(e)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	// Start server
	log.Printf("Server starting on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
