package main

import (
	"log"

	"saathi-backend/config"

	"saathi-backend/routes"
	tutorialrepository "saathi-backend/tutorial-repository"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Connect to MongoDB
	if err := config.ConnectDB(); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// if err := test_data.SeedTutorials(); err != nil {
	// 	log.Fatalf("Failed to seed tutorials: %v", err)
	// }

	defer config.DisconnectDB()

	// Initialize MongoDB repository
	tutorialrepository.InitRepository()

	// Create Fiber application
	//app := fiber.New()
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500 MB
	})

	routes.SetupRoutes(app)

	log.Println("Server running on http://localhost:8080")

	defer config.DisconnectDB()

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
