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

func (r *TutorialRepository) UploadVideo(
	ctx context.Context,
	bucketName string,
	fileName string,
	file io.Reader,
) error {

	err := r.gcsService.UploadVideo(
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
) ([]dto.TutorialResponseDTO, error) {

	filter := bson.M{
		"roles": bson.M{
			"$in": roles,
		},
		"is_active": true,
	}

	projection := bson.M{
		"video_id":          1,
		"app_name":          1,
		"video_title":       1,
		"video_description": 1,
		"created_at":        1,
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


func (r *TutorialRepository) GetUploadFiles(
	ctx context.Context,
	bucketName string,
	prefix string,
) (string, string, error) {

	videoPath, thumbnailPath, err := r.gcsService.GetUploadFiles(
		ctx,
		bucketName,
		prefix,
	)

	if err != nil {
		return "", "", fmt.Errorf(
			"failed to find upload files: %w",
			err,
		)
	}

	return videoPath, thumbnailPath, nil
}