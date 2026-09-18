package handlers

import (
	"context"
	"net/http"
	"time"

	"saathi-backend/dto"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/gofiber/fiber/v2"
)

type TutorialHandler struct {
	Service *tutorialservice.TutorialService
}

func NewTutorialHandler(
	service *tutorialservice.TutorialService,
) *TutorialHandler {
	return &TutorialHandler{
		Service: service,
	}
}

func (h *TutorialHandler) CreateTutorial(c *fiber.Ctx) error {

	var req dto.CreateTutorialRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(
		c.Context(),
		10*time.Second,
	)
	defer cancel()

	tutorial, err := h.Service.CreateTutorial(ctx, req)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Tutorial created successfully",
		"data":    tutorial,
	})
}

func (h *TutorialHandler) CreateTranslation(
	c *fiber.Ctx,
) error {

	// -----------------------------------------
	// 1. Parse request body
	// -----------------------------------------

	var req dto.CreateTranslationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// -----------------------------------------
	// 2. Create context with timeout
	// -----------------------------------------

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// -----------------------------------------
	// 3. Call service
	// -----------------------------------------

	translation, err := h.Service.CreateTranslation(
		ctx,
		req,
	)

	if err != nil {

		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{

			"success": false,
			"message": err.Error(),
		})
	}

	// -----------------------------------------
	// 4. Return response
	// -----------------------------------------

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{

		"success": true,
		"message": "Translation created successfully",
		"data":    translation,
	})
}
