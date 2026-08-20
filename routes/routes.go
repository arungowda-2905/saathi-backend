package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/saathi/api")

	api.Get("/v1/:id", handlers.GetTutorialByID)
	api.Get("/details/v1", handlers.GetDetailsByRole)

}
