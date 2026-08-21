package handlers

import (
	"context"
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

func GetDetailsByRole(c *fiber.Ctx) error {

	userRole := c.Get("X-Roles")

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
		"_id":               1,
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

// func GetAllTutorials(c *fiber.Ctx) error {

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	cursor, err := config.DB.Collection("tutorials").Find(ctx, bson.M{})
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to fetch tutorials",
// 		})
// 	}
// 	defer cursor.Close(ctx)

// 	var tutorials []model.Tutorial

// 	if err := cursor.All(ctx, &tutorials); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to decode tutorials",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(tutorials)
// }

// func GetTutorialsByAppName(c *fiber.Ctx) error {

// 	appName := c.Params("appName")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	cursor, err := config.DB.Collection("tutorials").Find(ctx, bson.M{"app_name": appName})
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to fetch tutorials",
// 		})
// 	}
// 	defer cursor.Close(ctx)

// 	var tutorials []model.Tutorial
// 	if err := cursor.All(ctx, &tutorials); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to decode tutorials",
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(tutorials)
// }
