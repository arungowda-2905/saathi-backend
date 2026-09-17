package tutorialrepository

import (
	"context"
	"fmt"

	"saathi-backend/config"
	"saathi-backend/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var collection *mongo.Collection

func InitRepository() {
	collection = config.DB.Collection("tutorial_translations")
}

func CreateTranslation(
	ctx context.Context,
	translation model.TutorialTranslation,
) error {

	_, err := collection.InsertOne(ctx, translation)
	fmt.Println("Inserting translation:", translation)
	return err
}
