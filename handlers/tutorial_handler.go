package handlers

import (
	"github.com/gofiber/fiber/v2"

	tutorialservice "saathi-backend/tutorial-service"
)

func CreateTutorialTranslation(c *fiber.Ctx) error {

	tutorialID := c.Params("tutorialID")
	language := c.FormValue("language")
	videoTitle := c.FormValue("video_title")
	videoDescription := c.FormValue("video_description")
	videoBucket := c.FormValue("video_bucket")
	titleImage := c.FormValue("title_image")
	duration := c.FormValue("duration")

	translation, err := tutorialservice.CreateTutorialTranslation(
		tutorialID,
		language,
		videoTitle,
		videoDescription,
		videoBucket,
		titleImage,
		duration,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tutorial translation created successfully",
		"data":    translation,
	})
}
