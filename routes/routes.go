package routes

import (
	"saathi-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, tutorialHandler *handlers.TutorialHandler) {

	api := app.Group("/saathi/api")

	api.Post("/v1/upload", tutorialHandler.UploadVideo)
	api.Post("/videos/v1", tutorialHandler.HandleVideoUpload)
	api.Patch("/v1/:videoId/translations", tutorialHandler.AddTutorialTranslation)
	api.Get("/v1/:videoId", tutorialHandler.GetTutorialByID)
	api.Get("/details/v1", tutorialHandler.GetDetailsByRole)
	api.Get("/v1/tutorials/app", tutorialHandler.GetDetailsByAppName)
	api.Post("/videos/:video_id/progress", handlers.CreateVideoProgress)
	api.Patch("/videos/:video_id/progress/feedback", handlers.UpdateVideoFeedback)
	api.Get("/v1/thumbnail/:videoId", tutorialHandler.GetTutorialThumbnailByID)
	api.Patch("/videos/:video_id/progress", handlers.UpdateVideoProgress)
	api.Get("/v1/videos/:video_id/progress", handlers.GetVideoProgress)
}
