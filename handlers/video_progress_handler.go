package handlers

import (
	"errors"
	"time"

	"saathi-backend/dto"
	"saathi-backend/model"
	"saathi-backend/tutorial-repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const hardcodedEmployeeID = "EMP001"

func CreateVideoProgress(c *fiber.Ctx) error {
	videoID := c.Params("video_id")

	if videoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "video_id is required",
		})
	}

	var req dto.CreateVideoProgressRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate position
	if req.PositionSeconds < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "position_seconds cannot be negative",
		})
	}

	// Check whether progress already exists
	existingProgress, err := tutorialrepository.GetVideoProgress(
		hardcodedEmployeeID,
		videoID,
	)

	if err == nil && existingProgress != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "video progress already exists",
		})
	}

	// If error is something other than "document not found",
	// return it instead of creating a new record.
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check existing video progress",
		})
	}

	now := time.Now()

	progress := &model.VideoProgress{
		ID:              uuid.New().String(),
		UserID:          hardcodedEmployeeID,
		VideoID:         videoID,
		PositionSeconds: req.PositionSeconds,
		Completed:       req.Completed,
		LastWatchedAt:   now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = tutorialrepository.CreateVideoProgress(progress)

	if err != nil {
		// Handles race condition where two POST requests
		// arrive at almost the same time.
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "video progress already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create video progress",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"video_id":         progress.VideoID,
		"position_seconds": progress.PositionSeconds,
		"completed":        progress.Completed,
		"last_watched_at":  progress.LastWatchedAt,
		"updated_at":       progress.UpdatedAt,
	})
}

func UpdateVideoFeedback(c *fiber.Ctx) error {
    videoID := c.Params("video_id")

    if videoID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "video_id is required",
        })
    }

    var req dto.VideoFeedbackRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    if req.Rating != "helpful" && req.Rating != "not_helpful" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "rating must be helpful or not_helpful",
        })
    }

    // Check whether progress already exists
    progress, err := tutorialrepository.GetVideoProgress(
        hardcodedEmployeeID,
        videoID,
    )

    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "video progress not found",
            })
        }

        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to check video progress",
        })
    }

    if progress == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "video progress not found",
        })
    }

    // Update only rating
    err = tutorialrepository.UpdateVideoFeedback(
        hardcodedEmployeeID,
        videoID,
        req.Rating,
    )

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to update feedback",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "video_id": videoID,
        "rating":   req.Rating,
    })
}