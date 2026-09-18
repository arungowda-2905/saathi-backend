package tutorialrepository

import (
	"context"
	"errors"
	"regexp"
	"saathi-backend/model"

	// "saathi-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// var collection *mongo.Collection
// var tutorialCollection *mongo.Collection

//	func InitRepository() {
//		collection = config.DB.Collection("tutorial_translations")
//	    tutorialCollection = config.DB.Collection("tutorials")
//	}
type TutorialRepository struct {
	Collection *mongo.Collection
}

func NewTutorialRepository(collection *mongo.Collection) *TutorialRepository {
	return &TutorialRepository{
		Collection: collection,
	}
}

func (r *TutorialRepository) CreateTutorial(
	ctx context.Context,
	tutorial *model.Tutorial,
) error {

	_, err := r.Collection.InsertOne(ctx, tutorial)

	return err
}

func (r *TutorialRepository) GetTutorialByTutorialID(
	ctx context.Context,
	tutorialID string,
) (*model.Tutorial, error) {

	var tutorial model.Tutorial

	err := r.Collection.FindOne(
		ctx,
		bson.M{
			"tutorial_id": tutorialID,
		},
	).Decode(&tutorial)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("tutorial not found")
		}

		return nil, err
	}

	return &tutorial, nil
}
func (r *TutorialRepository) AddLanguageToTutorial(
	ctx context.Context,
	tutorialID string,
	language string,
	translationID string,
) error {

	result, err := r.Collection.UpdateOne(
		ctx,
		bson.M{
			"tutorial_id": tutorialID,
		},
		bson.M{
			"$set": bson.M{
				"languages." + language: translationID,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("tutorial not found")
	}

	return nil
}

type TranslationRepository struct {
	Collection *mongo.Collection
}

func NewTranslationRepository(
	collection *mongo.Collection,
) *TranslationRepository {

	return &TranslationRepository{
		Collection: collection,
	}
}

func (r *TranslationRepository) CreateTranslation(
	ctx context.Context,
	translation *model.Translation,
) error {

	_, err := r.Collection.InsertOne(
		ctx,
		translation,
	)

	return err
}

//search

func (r *TutorialRepository) SearchTutorials(
	query string,
) ([]model.Tutorial, error) {

	ctx := context.Background()

	searchPattern := regexp.QuoteMeta(query)

	filter := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{
				"tutorial_title": bson.M{
					"$regex":   searchPattern,
					"$options": "i",
				},
			},
			{
				"tutorial_description": bson.M{
					"$regex":   searchPattern,
					"$options": "i",
				},
			},
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tutorials []model.Tutorial

	if err := cursor.All(ctx, &tutorials); err != nil {
		return nil, err
	}

	return tutorials, nil
}

func (r *TranslationRepository) GetTranslationsByIDs(
	translationIDs []string,
) (map[string]model.Translation, error) {

	ctx := context.Background()

	if len(translationIDs) == 0 {
		return make(map[string]model.Translation), nil
	}

	filter := bson.M{
		"translation_id": bson.M{
			"$in": translationIDs,
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var translations []model.Translation

	if err := cursor.All(ctx, &translations); err != nil {
		return nil, err
	}

	result := make(map[string]model.Translation)

	for _, translation := range translations {
		result[translation.TranslationID] = translation
	}

	return result, nil
}
