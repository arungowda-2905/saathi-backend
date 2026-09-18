package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	tutorialservice "saathi-backend/tutorial-service"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type videoService interface {
	GetVideo(videoID string, ctx context.Context) ([]byte, error)
	GetTutorialThumbnailByID(
		tid string,
		ctx context.Context,
	) ([]byte, string, error)
}

type TutorialHandler1 struct {
	service videoService
}

func NewTutorialHandler1(service videoService) *TutorialHandler1 {
	return &TutorialHandler1{service: service}
}

func (h *TutorialHandler1) GetVideo(c *fiber.Ctx) error {

	videoID := c.Params("videoId")

	if strings.TrimSpace(videoID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	videoBytes, err := h.service.GetVideo(videoID, ctx)

	if err != nil {
		log.Printf("ERROR in GetVideo: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//defer reader.Close()

	c.Set("Content-Type", "video/mp4")

	return c.Send(videoBytes)
}

func (h *TutorialHandler1) GetTutorialThumbnailByID(c *fiber.Ctx) error {

	tid := c.Params("thumbnailId")

	if tid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "thumbnail ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	imageBytes, contentType, err := h.service.GetTutorialThumbnailByID(
		tid,
		ctx,
	)

	if err != nil {

		fmt.Println("ERROR GetTutorialThumbnailByID:", err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "thumbnail not found",
			})
		}

		if errors.Is(err, tutorialservice.ErrInvalidThumbnailID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid thumbnail ID",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorial thumbnail",
		})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Length", strconv.Itoa(len(imageBytes)))
	c.Set("Cache-Control", "public, max-age=3600")

	return c.Send(imageBytes)
}

//search

func (h *TutorialHandler) SearchTutorials(c *fiber.Ctx) error {

	query := c.Query("query")
	language := c.Query("language")
	role := c.Query("X-Role")

	fmt.Println("QUERY:", query)
	fmt.Println("LANGUAGE:", language)
	fmt.Println("URL:", c.OriginalURL())

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "query is empty",
			},
		)
	}

	if language == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "language is empty",
			},
		)
	}

	if role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "role is empty",
			},
		)
	}

	results, err := h.Service.SearchTutorials(
		query,
		language,
		role,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    results,
		},
	)
}
