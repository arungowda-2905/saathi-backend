package routes

import (
	tutorialhandler "saathi-backend/tutorial-handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/saathi/api")

	// Upload video with metadata
	api.Post("/videos/v1", tutorialhandler.HandleVideoUpload)

	// Get tutorial by ID
	// api.Get("/tutorials/v1/:id", tutorialhandler.HandleGetTutorialByID)
}
