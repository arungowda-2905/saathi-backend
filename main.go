package main

import (
	"log"
	"saathi-backend/config"
	"saathi-backend/handlers"
	"saathi-backend/routes"
	tutorialrepository "saathi-backend/tutorial-repository"

	"github.com/gofiber/fiber/v2"
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
	app.Post("saathi/api/tutorials/v1", handlers.CreateTutorial)

	// Start server
	log.Println("Server running on http://localhost:8080")

	defer config.DisconnectDB()

	// err = test_data.SeedTutorials()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	routes.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
