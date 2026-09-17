package tutorialservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"saathi-backend/gcs"
)

type TutorialService struct {
	prefix string
}

func NewTutorialService(prefix string) *TutorialService {
	return &TutorialService{prefix: prefix}
}

var ErrInvalidVideoID = fmt.Errorf("invalid video ID")
var ErrInvalidThumbnailID = fmt.Errorf("invalid thumbnail ID")

func (s *TutorialService) GetVideo(
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

func (s *TutorialService) GetTutorialThumbnailByID(
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
