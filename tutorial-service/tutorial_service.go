package tutorialservice

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"
)

func CreateTutorialTranslation(
	tutorialID string,
	language string,
	videoTitle string,
	videoDescription string,
	videoBucket string,
	titleImage string,
	duration string,
) (*model.TutorialTranslation, error) {

	// Validate required fields
	if tutorialID == "" {
		return nil, errors.New("tutorial_id is required")
	}

	if language == "" {
		return nil, errors.New("language is required")
	}

	if videoTitle == "" {
		return nil, errors.New("video_title is required")
	}

	if videoDescription == "" {
		return nil, errors.New("video_description is required")
	}

	if videoBucket == "" {
		return nil, errors.New("video_bucket is required")
	}

	if titleImage == "" {
		return nil, errors.New("title_image is required")
	}

	if duration == "" {
		return nil, errors.New("duration is required")
	}

	// Generate timestamps
	now := time.Now()

	translation := &model.TutorialTranslation{
		ID:               primitive.NewObjectID(),
		TutorialID:       tutorialID,
		Language:         language,
		VideoTitle:       videoTitle,
		VideoDescription: videoDescription,
		VideoBucket:      videoBucket,
		TitleImage:       titleImage,
		Duration:         duration,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Save to MongoDB
	err := tutorialrepository.CreateTranslation(
		context.Background(),
		*translation,
	)

	if err != nil {
		return nil, err
	}

	return translation, nil
}
