package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"saathi-backend/config"
	tutorialhandler "saathi-backend/tutorial-handler"
	tutorialrepository "saathi-backend/tutorial-repository"
)

func main() {

	// Connect to MongoDB
	err := config.ConnectDB()
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Initialize repository
	tutorialrepository.InitRepository()

	// Create Fiber app
	app := fiber.New()

	// POST API
	app.Post("saathi/api/tutorials/v1", tutorialhandler.CreateTutorial)

	// Start server
	log.Println("Server running on http://localhost:8080")

	defer config.DisconnectDB()

	err = test_data.SeedTutorials()
	if err != nil {
		log.Fatal(err)
	}
	routes.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}