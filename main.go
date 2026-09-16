package main

import (
	"context"
	"log"

	"saathi-backend/config"
	"saathi-backend/gcs"

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

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://127.0.0.1:3000",
		AllowHeaders: "Origin, Content-Type, Accept, X-Role, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	//routes.SetupRoutes(app, tutorialHandler)

	log.Println("Server running on http://localhost:8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
