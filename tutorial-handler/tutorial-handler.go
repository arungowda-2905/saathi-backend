package tutorialhandler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

func HandleVideoUpload(c *fiber.Ctx) error {

	// Get video file
	videoFile, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video file is required",
		})
	}

	// Get metadata sent from frontend FormData
	title := c.FormValue("title")
	description := c.FormValue("description")
	applicationName := c.FormValue("applicationName")
	assignToRole := c.FormValue("assignToRole")

	// Validate empty values
	if title == "" ||
		description == "" ||
		applicationName == "" ||
		assignToRole == "" {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All metadata fields are required",
		})
	}

	// Create Tutorial object
	tutorial := model.Tutorial{
		AppName:          applicationName,
		VideoTitle:       title,
		VideoDescription: description,

		// Temporary: later this will be the GCS video URL
		VideoLink: videoFile.Filename,

		Roles: []string{
			assignToRole,
		},

		// Temporary/default values because frontend
		// is currently not sending these fields
		Version:    "1.0",
		TitleImage: "",
		IsActive:   true,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	createdTutorial, err := tutorialservice.CreateNewTutorial(
		ctx,
		tutorial,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tutorial",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Video uploaded successfully",
		"data":    createdTutorial,
	})
}