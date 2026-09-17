package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *handlers.TutorialHandler) {

	api := app.Group("/saathi/api")

	api.Post("/videos/v1", tutorialHandler.HandleVideoUpload)
	api.Post("/translations/v1/:TutorialID", handlers.CreateTutorialTranslation)

}
