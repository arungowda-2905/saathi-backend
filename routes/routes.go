package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *handlers.TutorialHandler) {

	api := app.Group("/saathi/api")

	api.Get("/v1/:videoId", tutorialHandler.GetVideo)
	api.Get("/v1/thumbnail/:thumbnailId", tutorialHandler.GetTutorialThumbnailByID)
}
