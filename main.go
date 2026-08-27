package main

import (
	"context"
	"log"

	"saathi-backend/config"
	"saathi-backend/gcs"
	"saathi-backend/handlers"
	"saathi-backend/routes"
	"saathi-backend/test_data"
	tutorialrepository "saathi-backend/tutorial-repository"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	if err := test_data.SeedTutorials(); err != nil {
		log.Fatalf("Failed to seed tutorials: %v", err)
	}

	tutorialrepository.InitRepository()

	ctx := context.Background()

	// Create GCS service
	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize GCS service: %v", err)
	}
	defer func() {
		if err := gcsService.Close(); err != nil {
			log.Printf("Failed to close GCS service: %v", err)
		}
	}()

	// Create repository
	repository := tutorialrepository.NewTutorialRepository(gcsService)

	// Create tutorial service
	service := tutorialservice.NewTutorialService(repository)

	// Create handler
	tutorialHandler := handlers.NewTutorialHandler(service)

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173",
		AllowHeaders: "Origin, Content-Type, Accept, X-Role, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	routes.SetupRoutes(app, tutorialHandler)

	log.Println("Server running on http://localhost:8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
