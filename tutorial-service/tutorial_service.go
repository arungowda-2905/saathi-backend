package tutorialservice

import (
	"context"
	"errors"
	"time"

	"saathi-backend/dto"
	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tutorialRepository interface {
	CreateTutorial(context.Context, model.Tutorial) error
	CreateTutorialTranslation(context.Context, model.TutorialTranslation) error
}

type TutorialService struct {
	repository tutorialRepository
}

func NewTutorialService(
	repository tutorialRepository,
) *TutorialService {
	return &TutorialService{
		repository: repository,
	}
}

func (s *TutorialService) HandleVideoUpload(
	ctx context.Context,
	req dto.HandleVideoUploadRequest,
) (fiber.Map, error) {

	now := time.Now()

	// Generate ONE unique UUID
	tutorialID := uuid.New().String()

	// 1. Create Tutorial parent

	tutorial := model.Tutorial{
		TutorialID: tutorialID,
		AppName:    req.AppName,
		Version:    req.Version,
		Roles:      req.Roles,
		IsActive:   req.IsActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.repository.CreateTutorial(ctx, tutorial)
	if err != nil {
		return nil, err
	}

	// 2. Create TutorialTranslation child

	translation := model.TutorialTranslation{

		TutorialID:       tutorialID,
		Language:         req.Language,
		VideoTitle:       req.VideoTitle,
		VideoDescription: req.VideoDescription,
		VideoBucket:      req.VideoBucket,
		TitleImage:       req.TitleImage,
		Duration:         req.Duration,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err = s.repository.CreateTutorialTranslation(ctx, translation)
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"message":     "Tutorial and translation created successfully",
		"tutorial":    tutorial,
		"translation": translation,
	}, nil
}

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
