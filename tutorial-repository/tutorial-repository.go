package tutorialrepository

import (
	"context"

	"saathi-backend/config"
	"saathi-backend/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var tutorialCollection *mongo.Collection

func InitRepository() {
	tutorialCollection = config.DB.Collection("tutorials")
}

func CreateTutorial(ctx context.Context, tutorial model.Tutorial) error {

	_, err := tutorialCollection.InsertOne(ctx, tutorial)

	return err
}