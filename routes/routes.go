package routes

import (
	handlers "saathi-backend/tutorial-handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/saathi/api")

	api.Get("/v1/:id", handlers.GetTutorialByID)

}
