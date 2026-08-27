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
	gcsFileName string,
) (model.Tutorial, error) {

	if gcsFileName == "" {
		return model.Tutorial{}, fmt.Errorf(
			"gcs file name is required",
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

	// Generate UUID for the video.
	videoID := uuid.New()

	// UUID is exposed to the application.
	tutorial.Video_ID = videoID

	// Actual GCS filename is stored internally.
	tutorial.Video_Bucket = gcsFileName

	// Set timestamps.
	now := time.Now()

	tutorial.CreatedAt = now
	tutorial.UpdatedAt = now

	// Save tutorial metadata.
	err := s.repository.InsertTutorial(
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
		//userRole,
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
	defer func() {
		if closeErr := gcsService.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close GCS service: %w", closeErr)
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

	return videoBytes, err
}

func (s *TutorialService) UploadVideo(
	ctx context.Context,
	originalFileName string,
	file io.Reader,
) (string, error) {

	if s.bucketName == "" {
		return "", fmt.Errorf("GCS_VIDEO_BUCKET is not configured")
	}

	extension := strings.ToLower(filepath.Ext(originalFileName))

	fileName := fmt.Sprintf(
		"%d%s",
		time.Now().UnixNano(),
		extension,
	)

	// Full path used for GCS
	objectName := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		fileName,
	)

	err := s.repository.UploadVideo(
		ctx,
		s.bucketName,
		objectName,
		file,
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload video: %w", err)
	}

	// Return only filename
	return fileName, nil
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
