package tutorialservice

import (
	"context"
	"time"

	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateTutorial(ctx context.Context, tutorial model.Tutorial) (model.Tutorial, error) {

	// Automatically generate unique video/tutorial ID
	tutorial.ID = bson.NewObjectID()

	// Automatically create timestamps
	now := time.Now()

	tutorial.CreatedAt = now
	tutorial.UpdatedAt = now

	// Save tutorial in MongoDB
	err := tutorialrepository.CreateTutorial(ctx, tutorial)
	if err != nil {
		return model.Tutorial{}, err
	}

	return tutorial, nil
}