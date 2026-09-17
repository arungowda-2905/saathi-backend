package routes

import (
	tutorialhandler "saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *tutorialhandler.TutorialHandler) {

	api := app.Group("/saathi/api")

	api.Post("/v1/tutorials", tutorialHandler.CreateTutorial)
	api.Post("/v1/translations", tutorialHandler.CreateTranslation)

}
