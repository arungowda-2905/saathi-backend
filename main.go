package main

import (
	"log"

	"saathi-backend/config"
	tutorialhandler "saathi-backend/handlers"
	"saathi-backend/routes"
	tutorialrepository "saathi-backend/tutorial-repository"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Connect to MongoDB
	if err := config.ConnectDB(); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Disconnect MongoDB when application stops
	defer func() {
		if err := config.DisconnectDB(); err != nil {
			log.Printf("MongoDB disconnect failed: %v", err)
		}
	}()

	// Tutorial repository
	tutorialRepository := tutorialrepository.NewTutorialRepository(
		config.DB.Collection("tutorials"),
	)

	// Translation repository
	translationRepository := tutorialrepository.NewTranslationRepository(
		config.DB.Collection("translations"),
	)

	// Tutorial service
	tutorialService := tutorialservice.NewTutorialService(
		tutorialRepository,
		translationRepository,
		//config.Client,
	)

	// Tutorial handler
	tutorialHandler := tutorialhandler.NewTutorialHandler(
		tutorialService,
	)

	// Fiber app
	app := fiber.New()

	// Routes
	routes.SetupRoutes(app, tutorialHandler)

	log.Println("Server running on http://localhost:8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
