package tutorialservice

import (
	"context"
	
	"time"

	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateNewTutorial(
	ctx context.Context,
	tutorial model.Tutorial,
) (model.Tutorial, error) {

	
	// Generate unique ID
	tutorial.ID = bson.NewObjectID()

	// Generate timestamps
	now := time.Now()
	tutorial.VideoLink = "Hardcoded"
	tutorial.CreatedAt = now
	tutorial.UpdatedAt = now

	// Insert into MongoDB
	err := tutorialrepository.InsertTutorial(ctx, tutorial)
	if err != nil {
		return model.Tutorial{}, err
	}

	return tutorial, nil
}