package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"saathi-backend/config"
	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

func CreateTutorial(c *fiber.Ctx) error {

	var tutorial model.Tutorial

	// Read JSON body
	if err := c.BodyParser(&tutorial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tutorial",
		})
	}

	// Return created tutorial
	return c.Status(fiber.StatusCreated).JSON(createdTutorial)
}

func GetTutorialByID(c *fiber.Ctx) error {

	id := c.Params("id")

	// Convert string ID to MongoDB ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tutorial ID",
		})
	}

	// Get user's role
	userRole := c.Get("X-Role")

	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role is required",
		})
	}

	collection := config.DB.Collection("tutorials")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tutorial model.Tutorial

	err = collection.FindOne(
		ctx,
		bson.M{
			"_id":       objectID,
			"is_active": true,
		},
	).Decode(&tutorial)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tutorial not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check whether user's role is allowed
	allowed := false

	for _, role := range tutorial.Roles {
		if role == userRole {
			allowed = true
			break
		}
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to access this tutorial",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tutorial)
}
