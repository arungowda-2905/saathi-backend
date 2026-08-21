package tutorialrepository

import (
	"context"

	"saathi-backend/config"
	"saathi-backend/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var tutorialCollection *mongo.Collection

func InitRepository() {
	tutorialCollection = config.DB.Collection("tutorials")
}

func InsertTutorial(
	ctx context.Context,
	tutorial model.Tutorial,
) error {

	_, err := tutorialCollection.InsertOne(ctx, tutorial)

	return err
}

func GetTutorialByID(
	ctx context.Context,
	videoID uuid.UUID,
	userRole string,
) (model.Tutorial, error) {

	var tutorial model.Tutorial

	filter := bson.M{
		"video_id":  videoID,
		"is_active": true,
		"roles": bson.M{
			"$in": []string{userRole},
		},
	}

	err := tutorialCollection.
		FindOne(ctx, filter).
		Decode(&tutorial)

	if err != nil {
		return model.Tutorial{}, err
	}

	return tutorial, nil
}
