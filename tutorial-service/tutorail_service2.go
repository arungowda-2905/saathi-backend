package tutorialservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"

	"saathi-backend/gcs"
)

type TutorialService1 struct {
	prefix string
}

func NewTutorialService1(prefix string) *TutorialService1 {
	return &TutorialService1{prefix: prefix}
}

var ErrInvalidVideoID = fmt.Errorf("invalid video ID")
var ErrInvalidThumbnailID = fmt.Errorf("invalid thumbnail ID")

func (s *TutorialService1) GetVideo(
	videoID string,
	ctx context.Context,
) ([]byte, error) {

	if strings.TrimSpace(videoID) == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	bucketName := os.Getenv("GCS_VIDEO_BUCKET")

	if strings.TrimSpace(bucketName) == "" {
		return nil, fmt.Errorf(
			"GCS_VIDEO_BUCKET is not configured",
		)
	}

	videoPath := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"video",
		videoID,
	)

	log.Printf("GCS bucket: %s", bucketName)
	log.Printf("GCS video path: %s", videoPath)

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
		gcsService.Close()

		return nil, fmt.Errorf(
			"failed to fetch video from GCS: %w",
			err,
		)
	}

	return videoBytes, nil
}

func (s *TutorialService1) GetTutorialThumbnailByID(
	thumbnailID string,
	ctx context.Context,
) ([]byte, string, error) {

	if strings.TrimSpace(thumbnailID) == "" {
		return nil, "", fmt.Errorf("thumbnail ID is required")
	}

	bucketName := os.Getenv("GCS_VIDEO_BUCKET")

	if strings.TrimSpace(bucketName) == "" {
		return nil, "", fmt.Errorf(
			"GCS_VIDEO_BUCKET is not configured",
		)
	}

	thumbnailPath := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"thumbnail",
		thumbnailID,
	)

	log.Printf("GCS bucket: %s", bucketName)
	log.Printf("GCS thumbnail path: %s", thumbnailPath)

	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to initialize GCS service: %w",
			err,
		)
	}

	imageBytes, contentType, err := gcsService.GetImage(
		ctx,
		bucketName,
		thumbnailPath,
	)

	if err != nil {
		gcsService.Close()

		return nil, "", fmt.Errorf(
			"failed to fetch thumbnail from GCS: %w",
			err,
		)
	}

	defer gcsService.Close()

	return imageBytes, contentType, nil
}

//for search

func (s *TutorialService) SearchTutorials(
	query string,
	selectedLanguage string,
	role string,
) ([]fiber.Map, error) {

	query = strings.TrimSpace(query)
	selectedLanguage = strings.TrimSpace(selectedLanguage)

	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if selectedLanguage == "" {
		return nil, fmt.Errorf("language cannot be empty")
	}

	// Search tutorials
	tutorials, err := s.Repository.SearchTutorials(query)
	if err != nil {
		return nil, err
	}

	if len(tutorials) == 0 {
		return []fiber.Map{}, nil
	}

	// Collect translation IDs
	translationIDs := make([]string, 0)

	for _, tutorial := range tutorials {

		translationID, exists :=
			tutorial.Languages[selectedLanguage]

		if exists && translationID != "" {
			translationIDs = append(
				translationIDs,
				translationID,
			)
		}
	}

	// Get translations
	translations, err :=
		s.TranslationRepository.GetTranslationsByIDs(
			translationIDs,
		)

	if err != nil {
		return nil, err
	}

	// Build response
	results := make([]fiber.Map, 0, len(tutorials))

	for _, tutorial := range tutorials {

		result := fiber.Map{
			"tutorialId":          tutorial.TutorialID,
			"appName":             tutorial.AppName,
			"tutorialTitle":       tutorial.TutorialTitle,
			"tutorialDescription": tutorial.TutorialDescription,
			"videoBucket":         tutorial.VideoBucket,
			"titleImage":          tutorial.ThumbnailImage,
			"duration":            tutorial.Duration,
			"language":            "English",
			"version":             tutorial.Version,
			"roles":               tutorial.Roles,
			"isActive":            tutorial.IsActive,
		}

		translationID, exists :=
			tutorial.Languages[selectedLanguage]

		if exists && translationID != "" {

			translation, found :=
				translations[translationID]

			if found {

				result["tutorialTitle"] =
					translation.TutorialTitle

				result["tutorialDescription"] =
					translation.TutorialDescription

				result["videoBucket"] =
					translation.VideoBucket

				result["titleImage"] =
					translation.ThumbnailImage

				result["duration"] =
					translation.Duration

				result["language"] =
					translation.Language
			}
		}

		results = append(results, result)
	}

	return results, nil
}
