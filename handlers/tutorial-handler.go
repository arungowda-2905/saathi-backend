package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"saathi-backend/config"
	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

func GetTutorialByID(c *fiber.Ctx) error {

	vid := c.Params("videoId")
	if vid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	userRole := c.Get("X-Role")

	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role is required",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	videoBytes, err := tutorialservice.GetTutorialByID(

		vid,
		ctx,
		userRole,
	)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tutorial not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorial",
		})
	}

	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Length", strconv.Itoa(len(videoBytes)))

	return c.Send(videoBytes)
}

func GetDetailsByRole(c *fiber.Ctx) error {

	userRole := c.Get("X-Role")

	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role not found",
		})
	}

	roles := strings.Split(userRole, ",")

	for i := range roles {
		roles[i] = strings.TrimSpace(roles[i])
	}

	filter := bson.M{
		"roles":     bson.M{"$in": roles},
		"is_active": true,
	}

	projection := bson.M{
		"video_id":          1,
		"app_name":          1,
		"video_title":       1,
		"video_description": 1,
	}

	opts := options.Find().SetProjection(projection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("tutorials").Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorials",
		})
	}
	defer cursor.Close(ctx)

	var tutorials []model.TutorialDetail

	if err := cursor.All(ctx, &tutorials); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tutorials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tutorials)
}

type TutorialHandler struct {
	service *tutorialservice.TutorialService
}

func NewTutorialHandler(service *tutorialservice.TutorialService) *TutorialHandler {
	return &TutorialHandler{
		service: service,
	}
}

func (h *TutorialHandler) UploadVideo(c *fiber.Ctx) error {
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video file is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded video",
		})
	}
	defer src.Close()

	fileName, err := h.service.UploadVideo(
		c.Context(),
		file.Filename,
		src,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload video",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Video uploaded successfully",
		"fileName": fileName,
	})
}
