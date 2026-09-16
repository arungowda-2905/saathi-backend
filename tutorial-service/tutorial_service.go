package tutorialservice

import (
	"context"
	"time"

	"saathi-backend/dto"
	"saathi-backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	// Same ID is used as VideoID
	videoID := tutorialID

	// 1. Create Tutorial parent

	tutorial := model.Tutorial{
		ID:        tutorialID,
		VideoID:   videoID,
		AppName:   req.AppName,
		Version:   req.Version,
		Roles:     req.Roles,
		IsActive:  req.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repository.CreateTutorial(ctx, tutorial)
	if err != nil {
		return nil, err
	}

	// 2. Create TutorialTranslation child

	translation := model.TutorialTranslation{
		ID:               uuid.New().String(),
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
