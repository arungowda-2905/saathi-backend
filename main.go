package main

import (
	"context"
	"log"

	"saathi-backend/config"
	"saathi-backend/gcs"
	"saathi-backend/handlers"
	"saathi-backend/routes"
	tutorialrepository "saathi-backend/tutorial-repository"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal("Error loading dev.env")
	}

	if err := config.ConnectDB(); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	defer config.DisconnectDB()

	// if err := test_data.SeedTutorials(); err != nil {
	// 	log.Fatalf("Failed to seed tutorials: %v", err)
	// }

	tutorialrepository.InitRepository()

	ctx := context.Background()

	// Create GCS service
	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize GCS service: %v", err)
	}

	// Create repository
	repository := tutorialrepository.NewTutorialRepository(gcsService)

	// Create tutorial service
	service := tutorialservice.NewTutorialService(repository)

	// Create handler
	tutorialHandler := handlers.NewTutorialHandler(service)

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	routes.SetupRoutes(app, tutorialHandler)

	log.Println("Server running on http://localhost:8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
