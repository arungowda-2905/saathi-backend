package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *handlers.TutorialHandler) {

	api := app.Group("/saathi/api")

	api.Post("/v1/upload", tutorialHandler.UploadVideo)
	api.Post("/videos/v1", handlers.HandleVideoUpload)

	api.Get("/v1/:videoId", handlers.GetTutorialByID)
	api.Get("/details/v1", handlers.GetDetailsByRole)
}
