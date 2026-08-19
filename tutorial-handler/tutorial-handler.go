package tutorialhandler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"

	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

func CreateTutorial(c fiber.Ctx) error {

	var tutorial model.Tutorial

	// Read JSON body
	if err := c.Bind().Body(&tutorial); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Call service
	createdTutorial, err := tutorialservice.CreateTutorial(ctx, tutorial)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create tutorial",
		})
	}

	// Return created tutorial
	return c.Status(201).JSON(createdTutorial)
}