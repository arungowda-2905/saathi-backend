package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *handlers.TutorialHandler, tutorialHandler1 *handlers.TutorialHandler1) {

	api := app.Group("/saathi/api")

	api.Post("/v1/upload", tutorialHandler.UploadVideo)
	api.Post("/v1/tutorials", tutorialHandler.CreateTutorial)
	api.Post("/v1/translations", tutorialHandler.CreateTranslation)
	api.Get("/v1/details", tutorialHandler.GetDetailsByRole)
	api.Get("/v1/:videoId", tutorialHandler1.GetVideo)
	api.Get("/v1/thumbnail/:thumbnailId", tutorialHandler1.GetTutorialThumbnailByID)
}
