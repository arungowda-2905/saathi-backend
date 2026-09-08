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

var ErrAppNameNotFound = errors.New("app name not found")
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
	videoBucketUUID string,
	thumbnailUUID string,
) (model.Tutorial, error) {

	if videoBucketUUID == "" || thumbnailUUID == "" {
		return model.Tutorial{}, fmt.Errorf(
			"video bucket UUID and thumbnail UUID are required",
		)
	}

	if tutorial.AppName == "" {
		return model.Tutorial{}, fmt.Errorf(
			"application name is required",
		)
	}

	if len(tutorial.Roles) == 0 {
		return model.Tutorial{}, fmt.Errorf(
			"at least one role is required",
		)
	}

	if len(tutorial.Translations) == 0 {
		return model.Tutorial{}, fmt.Errorf("at least one translation is required")
	}

	for language, translation := range tutorial.Translations {
		if translation.VideoTitle == "" || translation.VideoDescription == "" {
			return model.Tutorial{}, fmt.Errorf("title and description are required for language %s", language)
		}
		if translation.Duration == "" {
			return model.Tutorial{}, fmt.Errorf("duration is required for language %s", language)
		}
	}

	videoPrefix := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"video",
		videoBucketUUID,
	)
	thumbnailPrefix := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"thumbnail",
		thumbnailUUID,
	)

	videoPath, err := s.repository.GetUploadFile(
		ctx,
		s.bucketName,
		videoPrefix,
		"video",
	)
	if err != nil {
		return model.Tutorial{}, fmt.Errorf(
			"failed to find uploaded video: %w",
			err,
		)
	}

	thumbnailPath, err := s.repository.GetUploadFile(
		ctx,
		s.bucketName,
		thumbnailPrefix,
		"thumbnail",
	)

	if err != nil {
		return model.Tutorial{}, fmt.Errorf(
			"failed to find uploaded files: %w",
			err,
		)
	}

	// Generate final UUID for the tutorial.
	videoID := uuid.New()

	tutorial.Video_ID = videoID.String()

	for language, translation := range tutorial.Translations {
		translation.Video_Bucket = path.Base(videoPath)
		translation.TitleImage = path.Base(thumbnailPath)
		tutorial.Translations[language] = translation
	}

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

	if _, err := uuid.Parse(videoId); err != nil {
		return nil, fmt.Errorf(
			"%w: %s",
			ErrInvalidVideoID,
			videoId,
		)
	}

	tutorial, err := s.repository.GetTutorialByID(
		ctx,
		videoId,
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

	translation, ok := firstTranslation(tutorial.Translations)
	if !ok {
		return nil, fmt.Errorf("tutorial translation is missing")
	}

	videoPath := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"video",
		translation.Video_Bucket,
	)

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

func firstTranslation(translations map[string]model.Translation) (model.Translation, bool) {
	for _, translation := range translations {
		return translation, true
	}
	return model.Translation{}, false
}

func (s *TutorialService) UploadVideoAndThumbnail(
	ctx context.Context,
	videoFileName string,
	videoFile io.Reader,
	thumbnailFileName string,
	thumbnailFile io.Reader,
) (string, string, error) {

	if s.bucketName == "" {
		return "", "", fmt.Errorf("GCS_VIDEO_BUCKET is not configured")
	}

	// videoBucketUUID := uuid.New().String()
	// thumbnailUUID := uuid.New().String()
	videoBucketUnique := fmt.Sprintf("%d%s", time.Now().UnixNano(), strings.ToLower(filepath.Ext(videoFileName)))
	thumbnailUnique := fmt.Sprintf("%d%s", time.Now().UnixNano(), strings.ToLower(filepath.Ext(thumbnailFileName)))

	videoFolder := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"video",
	)
	thumbnailFolder := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"thumbnail",
	)

	// GCS object names
	videoObjectName := path.Join(
		videoFolder,
		videoBucketUnique,
	)

	thumbnailObjectName := path.Join(
		thumbnailFolder,
		thumbnailUnique,
	)

	// Upload video
	err := s.repository.UploadVideoThumbnail(
		ctx,
		s.bucketName,
		videoObjectName,
		videoFile,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to upload video: %w",
			err,
		)
	}

	// Upload thumbnail
	err = s.repository.UploadVideoThumbnail(
		ctx,
		s.bucketName,
		thumbnailObjectName,
		thumbnailFile,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to upload thumbnail: %w",
			err,
		)
	}

	return videoBucketUnique, thumbnailUnique, nil
}
func (s *TutorialService) GetDetailsByRole(
	ctx context.Context,
	userRole string,
) ([]dto.AppCountResponseDTO, error) {

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

func (s *TutorialService) GetDetailsByAppName(
	ctx context.Context,
	appName string,
	language string,
) ([]dto.TutorialResponseDTO, error) {

	appName = strings.TrimSpace(appName)
	language = strings.TrimSpace(strings.ToLower(language))

	if appName == "" {
		return nil, ErrAppNameNotFound
	}

	if language == "" {
		language = "en"
	}

	return s.repository.GetTutorialsByAppName(
		ctx,
		appName,
		language,
	)
}

func (s *TutorialService) AddTranslation(
	ctx context.Context,
	videoID string,
	language string,
	translation model.Translation,
) error {
	if _, err := uuid.Parse(videoID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidVideoID, videoID)
	}
	return s.repository.UpdateTranslation(ctx, videoID, language, translation)
}

func (s *TutorialService) GetTutorialThumbnailByID(
	videoID string,
	ctx context.Context,
) ([]byte, string, error) {

	if videoID == "" {
		return nil, "", fmt.Errorf("video ID is required")
	}

	if _, err := uuid.Parse(videoID); err != nil {
		return nil, "", fmt.Errorf(
			"%w: %s",
			ErrInvalidVideoID,
			videoID,
		)
	}

	tutorial, err := s.repository.GetTutorialByID(
		ctx,
		videoID,
	)

	if err != nil {
		return nil, "", err
	}

	bucketName := os.Getenv("GCS_VIDEO_BUCKET")

	if bucketName == "" {
		return nil, "", fmt.Errorf(
			"GCS_VIDEO_BUCKET is not configured",
		)
	}

	translation, ok := firstTranslation(tutorial.Translations)
	fmt.Println("arun", translation)
	if !ok {
		return nil, "", fmt.Errorf("tutorial translation is missing")
	}

	imagePath := path.Join(
		strings.TrimSuffix(s.prefix, "/"),
		"thumbnail",
		translation.TitleImage,
	)

	if imagePath == "" {
		return nil, "", fmt.Errorf("thumbnail path is empty")
	}

	gcsService, err := gcs.NewService(ctx)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to initialize GCS service: %w",
			err,
		)
	}

	defer func() {
		_ = gcsService.Close()
	}()

	imageBytes, contentType, err := gcsService.GetImage(
		ctx,
		bucketName,
		imagePath,
	)

	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to fetch thumbnail from GCS: %w",
			err,
		)
	}

	return imageBytes, contentType, nil
}
