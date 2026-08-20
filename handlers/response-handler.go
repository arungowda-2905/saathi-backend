package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"saathi-backend/config"
	"saathi-backend/model"
)

func GetDetailsByRole(c *fiber.Ctx) error {

	// Get roles from authentication middlewares
	// roles := c.Locals("roles")
	userRole := c.Get("X-Role")

	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role not found",
		})
	}

	userRoles := []string{userRole}

	if len(userRoles) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User has no roles",
		})
	}

	// Find active tutorials assigned to any
	// of the user's roles
	filter := bson.M{
		"roles": bson.M{
			"$in": userRoles,
		},
		"is_active": true,
	}

	// Return only required fields
	projection := bson.M{
		"_id":         0,
		"video_title": 1,
		"app_name":    1,
		"video_link":  1,
		"created_at":  1,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	cursor, err := config.DB.
		Collection("tutorials").
		Find(
			ctx,
			filter,
			options.Find().SetProjection(projection),
		)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorials",
		})
	}

	defer cursor.Close(ctx)

	var tutorials []model.TutorialRoleResponse

	if err := cursor.All(ctx, &tutorials); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tutorials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tutorials,
	})
}
