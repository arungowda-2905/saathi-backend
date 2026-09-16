package handlers

import (
	"context"

	"saathi-backend/dto"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/gofiber/fiber/v2"
)

type TutorialHandler struct {
	service *tutorialservice.TutorialService
}

func NewTutorialHandler(
	service *tutorialservice.TutorialService,
) *TutorialHandler {
	return &TutorialHandler{
		service: service,
	}
}

func (h *TutorialHandler) HandleVideoUpload(c *fiber.Ctx) error {

	var req dto.HandleVideoUploadRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AppName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "appName is required",
		})
	}

	if req.Version == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "version is required",
		})
	}

	if req.Language == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "language is required",
		})
	}

	ctx := context.Background()

	result, err := h.service.HandleVideoUpload(ctx, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
