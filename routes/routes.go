package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/saathi/api")
	api.Post("/translations/v1/:TutorialID", handlers.CreateTutorialTranslation)
}

