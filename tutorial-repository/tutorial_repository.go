package tutorialrepository

import (
	"context"

	"saathi-backend/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type TutorialRepository struct {
	tutorialCollection    *mongo.Collection
	translationCollection *mongo.Collection
}

var tutorialCollection *mongo.Collection

func NewTutorialRepository(db *mongo.Database) *TutorialRepository {
	return &TutorialRepository{
		tutorialCollection:    db.Collection("tutorials"),
		translationCollection: db.Collection("tutorial_translations"),
	}
}

func (r *TutorialRepository) CreateTutorial(
	ctx context.Context,
	tutorial model.Tutorial,
) error {

	_, err := r.tutorialCollection.InsertOne(ctx, tutorial)

	return err
}

func (r *TutorialRepository) CreateTutorialTranslation(
	ctx context.Context,
	translation model.TutorialTranslation,
) error {

	_, err := r.translationCollection.InsertOne(ctx, translation)

	return err
}

func CreateTranslation(
	ctx context.Context,
	translation model.TutorialTranslation,
) error {

	_, err := tutorialCollection.InsertOne(ctx, translation)

	return err
}
