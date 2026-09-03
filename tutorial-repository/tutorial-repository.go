package tutorialrepository

import (
	"context"
	"fmt"
	"io"

	"saathi-backend/config"
	"saathi-backend/dto"
	"saathi-backend/gcs"
	"saathi-backend/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var tutorialCollection *mongo.Collection

func InitRepository() {
	tutorialCollection = config.DB.Collection("tutorials")
	videoProgressCollection = config.DB.Collection("video_progress")
}

type TutorialRepository struct {
	gcsService *gcs.Service
}

func NewTutorialRepository(gcsService *gcs.Service) *TutorialRepository {
	return &TutorialRepository{
		gcsService: gcsService,
	}
}

func (r *TutorialRepository) InsertTutorial(
	ctx context.Context,
	tutorial model.Tutorial,
) error {

	_, err := tutorialCollection.InsertOne(
		ctx,
		tutorial,
	)

	return err
}

func (r *TutorialRepository) GetTutorialByID(
	ctx context.Context,
	videoID uuid.UUID,
	// userRole string,
) (model.Tutorial, error) {

	var tutorial model.Tutorial

	filter := bson.M{
		"video_id":  videoID,
		"is_active": true,
		// "roles": bson.M{
		// 	"$in": []string{userRole},
		// },
	}

	err := tutorialCollection.
		FindOne(ctx, filter).
		Decode(&tutorial)

	if err != nil {
		return model.Tutorial{}, err
	}

	return tutorial, nil
}

func (r *TutorialRepository) UploadVideoThumbnail(
	ctx context.Context,
	bucketName string,
	fileName string,
	file io.Reader,
) error {

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

	cursor, err := tutorialCollection.Aggregate(ctx, pipeline)
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

func (r *TutorialRepository) GetTutorialsByAppName(
	ctx context.Context,
	appName string,
) ([]dto.TutorialResponseDTO, error) {

	filter := bson.M{
		"app_name":  appName,
		"is_active": true,
	}

	projection := bson.M{
		"video_id":          1,
		"video_title":       1,
		"video_description": 1,
		"title_image":       1,
		"duration":          1,
	}

	opts := options.Find().SetProjection(projection)

	cursor, err := tutorialCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	tutorials := make([]dto.TutorialResponseDTO, 0)

	if err := cursor.All(ctx, &tutorials); err != nil {
		return nil, err
	}

	return tutorials, nil
}

func (r *TutorialRepository) GetUploadFile(
	ctx context.Context,
	bucketName string,
	prefix string,
	fileType string,
) (string, error) {
	filePath, err := r.gcsService.GetUploadFile(ctx, bucketName, prefix, fileType)
	if err != nil {
		return "", fmt.Errorf("failed to find upload file: %w", err)
	}

	return filePath, nil
}
