package tutorialservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"saathi-backend/dto"
	gcs "saathi-backend/gcs"
	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

var ErrRoleNotFound = errors.New("user role not found")
var ErrInvalidVideoID = fmt.Errorf("invalid video ID")

type TutorialService struct {
	repository *tutorialrepository.TutorialRepository
	bucketName string
	prefix     string
}

func NewTutorialService(
	repository *tutorialrepository.TutorialRepository,
) *TutorialService {
	return &TutorialService{
		repository: repository,
		bucketName: os.Getenv("GCS_VIDEO_BUCKET"),
		prefix:     os.Getenv("GCS_VIDEO_PREFIX"),
	}
}

func (s *TutorialService) CreateNewTutorial(
	ctx context.Context,
	tutorial model.Tutorial,
	uploadID string,
) (model.Tutorial, error) {

	if uploadID == "" {
		return model.Tutorial{}, fmt.Errorf(
			"upload ID is required",
		)
	}

	if tutorial.AppName == "" {
		return model.Tutorial{}, fmt.Errorf(
			"application name is required",
		)
	}

	if tutorial.VideoTitle == "" {
		return model.Tutorial{}, fmt.Errorf(
			"video title is required",
		)
	}

	if tutorial.VideoDescription == "" {
		return model.Tutorial{}, fmt.Errorf(
			"video description is required",
		)
	}

	if len(tutorial.Roles) == 0 {
		return model.Tutorial{}, fmt.Errorf(
			"at least one role is required",
		)
	}

	// Find video + thumbnail using upload_id.
	uploadPrefix := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"uploads",
		uploadID,
	)

	videoPath, thumbnailPath, err := s.repository.GetUploadFiles(
		ctx,
		s.bucketName,
		uploadPrefix,
	)

	if err != nil {
		return model.Tutorial{}, fmt.Errorf(
			"failed to find uploaded files: %w",
			err,
		)
	}

	// Generate final UUID for the tutorial.
	videoID := uuid.New()

	tutorial.Video_ID = videoID

	// Store internal GCS paths.
	tutorial.Video_Bucket = videoPath
	tutorial.TitleImage = thumbnailPath

	// Set timestamps.
	now := time.Now()

	tutorial.CreatedAt = now
	tutorial.UpdatedAt = now

	// Save metadata.
	err = s.repository.InsertTutorial(
		ctx,
		tutorial,
	)

	if err != nil {
		return model.Tutorial{}, fmt.Errorf(
			"failed to insert tutorial: %w",
			err,
		)
	}

	return tutorial, nil
}

func (s *TutorialService) GetTutorialByID(
	videoId string,
	ctx context.Context,
	// userRole string,
) (videoBytes []byte, err error) {

	if videoId == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	videoUUID, err := uuid.Parse(videoId)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %s",
			ErrInvalidVideoID,
			videoId,
		)
	}

	tutorial, err := s.repository.GetTutorialByID(
		ctx,
		videoUUID,
		// userRole,
	)
	if err != nil {
		return nil, err
	}

	bucketName := os.Getenv("GCS_VIDEO_BUCKET")

	if bucketName == "" {
		return nil, fmt.Errorf(
			"GCS_VIDEO_BUCKET is not configured",
		)
	}

	// Video_Bucket already contains the complete GCS object path.
	// Example:
	// uploads/<upload_id>/video.mp4
	videoPath := tutorial.Video_Bucket

	if videoPath == "" {
		return nil, fmt.Errorf("video path is empty")
	}

	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize GCS service: %w",
			err,
		)
	}

	defer func() {
		if closeErr := gcsService.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf(
				"failed to close GCS service: %w",
				closeErr,
			)
		}
	}()

	videoBytes, err = gcsService.GetVideo(
		ctx,
		bucketName,
		videoPath,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch video from GCS: %w",
			err,
		)
	}

	return videoBytes, nil
}
func (s *TutorialService) UploadVideoAndThumbnail(
	ctx context.Context,
	videoFileName string,
	videoFile io.Reader,
	thumbnailFileName string,
	thumbnailFile io.Reader,
) (string, error) {

	if s.bucketName == "" {
		return "", fmt.Errorf("GCS_VIDEO_BUCKET is not configured")
	}

	// Generate ONE UUID for both files.
	uploadID := uuid.New().String()

	// Create one GCS folder/prefix using the UUID.
	uploadPrefix := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"uploads",
		uploadID,
	)

	// Keep the original extensions internally.
	videoExtension := strings.ToLower(
		filepath.Ext(videoFileName),
	)

	thumbnailExtension := strings.ToLower(
		filepath.Ext(thumbnailFileName),
	)

	// GCS object names
	videoObjectName := path.Join(
		uploadPrefix,
		"video"+videoExtension,
	)

	thumbnailObjectName := path.Join(
		uploadPrefix,
		"thumbnail"+thumbnailExtension,
	)

	// Upload video
	err := s.repository.UploadVideo(
		ctx,
		s.bucketName,
		videoObjectName,
		videoFile,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to upload video: %w",
			err,
		)
	}

	// Upload thumbnail
	err = s.repository.UploadVideo(
		ctx,
		s.bucketName,
		thumbnailObjectName,
		thumbnailFile,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to upload thumbnail: %w",
			err,
		)
	}

	// Return only the UUID.
	// Actual GCS object names are never exposed.
	return uploadID, nil
}
func (s *TutorialService) GetDetailsByRole(
	ctx context.Context,
	userRole string,
) ([]dto.TutorialResponseDTO, error) {

	if strings.TrimSpace(userRole) == "" {
		return nil, ErrRoleNotFound
	}

	roles := strings.Split(userRole, ",")

	var cleanedRoles []string

	for _, role := range roles {
		role = strings.TrimSpace(role)

		if role != "" {
			cleanedRoles = append(cleanedRoles, role)
		}
	}

	if len(cleanedRoles) == 0 {
		return nil, ErrRoleNotFound
	}

	return s.repository.GetTutorialsByRoles(ctx, cleanedRoles)
}
