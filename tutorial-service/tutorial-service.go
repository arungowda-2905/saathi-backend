package tutorialservice

import (
	"context"
	"fmt"
	"os"
	"path"
	gcs "saathi-backend/gcs"
	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"
	"time"

	uuid "github.com/google/uuid"

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
	//tutorial.Video_ID = "Hardcoded"
	tutorial.CreatedAt = now
	tutorial.UpdatedAt = now

	// Insert into MongoDB
	err := tutorialrepository.InsertTutorial(ctx, tutorial)
	if err != nil {
		return model.Tutorial{}, err
	}

	return tutorial, nil
}

func GetTutorialByID(
	videoId string,
	ctx context.Context,
	userRole string,
) ([]byte, error) {

	if videoId == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	// Convert URL string to UUID
	videoUUID, err := uuid.Parse(videoId)
	if err != nil {
		return nil, fmt.Errorf("invalid video ID %q: %w", videoId, err)
	}

	// Find tutorial in MongoDB
	tutorial, err := tutorialrepository.GetTutorialByID(
		ctx,
		videoUUID,
		userRole,
	)
	if err != nil {
		return nil, err
	}

	bucketName := os.Getenv("GCS_VIDEO_BUCKET")
	videoPrefix := os.Getenv("GCS_VIDEO_PREFIX")

	if bucketName == "" {
		return nil, fmt.Errorf(
			"GCS_VIDEO_BUCKET is not configured",
		)
	}

	if videoPrefix == "" {
		return nil, fmt.Errorf(
			"GCS_VIDEO_PREFIX is not configured",
		)
	}

	videoFileName := tutorial.Video_Bucket

	videoPath := path.Join(
		videoPrefix,
		videoFileName,
	)

	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize GCS service: %w",
			err,
		)
	}

	videoBytes, err := gcsService.GetVideo(
		ctx,
		bucketName,
		videoPath,
	)

	if err != nil {
		return nil, err
	}

	return videoBytes, nil
}
