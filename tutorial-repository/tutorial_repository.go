package tutorialrepository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"saathi-backend/dto"
	"saathi-backend/gcs"
	"saathi-backend/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TutorialRepository struct {
	Collection *mongo.Collection
	gcsService *gcs.Service
}

func NewTutorialRepository(
	collection *mongo.Collection,
	gcsService *gcs.Service,
) *TutorialRepository {
	return &TutorialRepository{
		Collection: collection,
		gcsService: gcsService,
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

func (r *TutorialRepository) UploadVideoThumbnail(
	ctx context.Context,
	bucketName string,
	fileName string,
	file io.Reader,
) error {
	if r == nil || r.gcsService == nil {
		return fmt.Errorf("GCS service is not initialized")
	}

	err := r.gcsService.UploadVideoThumbnail(
		ctx,
		bucketName,
		fileName,
		file,
	)
	if err != nil {
		return fmt.Errorf("GCS upload error: %w", err)
	}

	return nil
}
func (r *TutorialRepository) GetTutorialsByRoles(
	ctx context.Context,
	roles []string,
) ([]dto.AppCountResponseDTO, error) {

	pipeline := mongo.Pipeline{
		// 1. Filter by role and active tutorials
		{
			{Key: "$match", Value: bson.D{
				{Key: "roles", Value: bson.D{
					{Key: "$in", Value: roles},
				}},
				{Key: "is_active", Value: true},
			}},
		},

		// 2. Group by app_name and count
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$app_name"},
				{Key: "count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
			}},
		},

		// 3. Convert _id to app_name
		{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "app_name", Value: "$_id"},
				{Key: "count", Value: 1},
			}},
		},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	appCounts := make([]dto.AppCountResponseDTO, 0)

	if err := cursor.All(ctx, &appCounts); err != nil {
		return nil, err
	}

	return appCounts, nil
}
