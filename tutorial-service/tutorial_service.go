package tutorialservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"saathi-backend/dto"
	"saathi-backend/model"
	tutorialrepository "saathi-backend/tutorial-repository"

	"github.com/google/uuid"
)

var ErrRoleNotFound = errors.New("user role not found")

type TutorialService struct {
	Repository            *tutorialrepository.TutorialRepository
	TranslationRepository *tutorialrepository.TranslationRepository
	// DB                    *mongo.Client
	bucketName string
	prefix     string
}

func NewTutorialService(
	repository *tutorialrepository.TutorialRepository,
	translationRepository *tutorialrepository.TranslationRepository,
	//client *mongo.Client,
) *TutorialService {
	bucketName := strings.TrimSpace(os.Getenv("GCS_VIDEO_BUCKET"))
	prefix := strings.TrimSpace(os.Getenv("GCS_VIDEO_PREFIX"))

	return &TutorialService{
		Repository:            repository,
		TranslationRepository: translationRepository,
		bucketName:            bucketName,
		prefix:                prefix,
		//DB:                    client,
	}
}

func (s *TutorialService) CreateTutorial(
	ctx context.Context,
	req dto.CreateTutorialRequest,
) (string, error) {

	// Basic validation
	if strings.TrimSpace(req.AppName) == "" {
		return "", errors.New("appName is required")
	}

	if strings.TrimSpace(req.TutorialTitle) == "" {
		return "", errors.New("tutorialTitle is required")
	}

	if strings.TrimSpace(req.TutorialDescription) == "" {
		return "", errors.New("tutorialDescription is required")
	}

	if strings.TrimSpace(req.VideoBucket) == "" {
		return "", errors.New("videoBucket is required")
	}

	if strings.TrimSpace(req.ThumbnailImage) == "" {
		return "", errors.New("titleImage is required")
	}

	if strings.TrimSpace(req.Duration) == "" {
		return "", errors.New("duration is required")
	}

	if strings.TrimSpace(req.Version) == "" {
		return "", errors.New("version is required")
	}

	if len(req.Roles) == 0 {
		return "", errors.New("at least one role is required")
	}

	// Generate unique UUID
	tutorialID := uuid.New().String()

	now := time.Now()

	tutorial := &model.Tutorial{
		TutorialID:          tutorialID,
		AppName:             req.AppName,
		TutorialTitle:       req.TutorialTitle,
		TutorialDescription: req.TutorialDescription,
		VideoBucket:         req.VideoBucket,
		ThumbnailImage:      req.ThumbnailImage,
		Duration:            req.Duration,

		// Languages MUST be empty during tutorial creation.
		Languages: make(map[string]string),

		Version:   req.Version,
		Roles:     req.Roles,
		IsActive:  req.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.Repository.CreateTutorial(ctx, tutorial)
	if err != nil {
		return "", err
	}

	return tutorial.TutorialID, nil
}

func (s *TutorialService) CreateTranslation(
	ctx context.Context,
	req dto.CreateTranslationRequest,
) (*model.Translation, error) {

	// 1. Validate Tutorial ID
	if strings.TrimSpace(req.TutorialID) == "" {
		return nil, errors.New("tutorialId is required")
	}

	// 2. Validate Language
	if strings.TrimSpace(req.Language) == "" {
		return nil, errors.New("language is required")
	}

	// 3. Validate Tutorial Title
	if strings.TrimSpace(req.TutorialTitle) == "" {
		return nil, errors.New("tutorialTitle is required")
	}

	// 4. Validate Description
	if strings.TrimSpace(req.TutorialDescription) == "" {
		return nil, errors.New("tutorialDescription is required")
	}

	// 5. Validate Video Bucket
	if strings.TrimSpace(req.VideoBucket) == "" {
		return nil, errors.New("videoBucket is required")
	}

	// 6. Validate Thumbnail
	if strings.TrimSpace(req.ThumbnailImage) == "" {
		return nil, errors.New("titleImage is required")
	}

	// 7. Validate Duration
	if strings.TrimSpace(req.Duration) == "" {
		return nil, errors.New("duration is required")
	}

	// 8. Check whether tutorial exists
	tutorial, err := s.Repository.GetTutorialByTutorialID(
		ctx,
		req.TutorialID,
	)
	if err != nil {
		return nil, err
	}

	// 9. Check whether language already exists
	if tutorial.Languages != nil {

		_, exists := tutorial.Languages[req.Language]

		if exists {
			return nil, errors.New(
				"translation for this language already exists",
			)
		}
	}

	// 10. Generate Translation UUID
	translationID := uuid.New().String()

	// 11. Create Translation object
	translation := &model.Translation{
		TranslationID:       translationID,
		TutorialTitle:       req.TutorialTitle,
		TutorialDescription: req.TutorialDescription,
		VideoBucket:         req.VideoBucket,
		ThumbnailImage:      req.ThumbnailImage,
		Duration:            req.Duration,
		Language:            req.Language,
	}

	// 12. Insert translation
	err = s.TranslationRepository.CreateTranslation(
		ctx,
		translation,
	)
	if err != nil {
		return nil, err
	}

	// 13. Update tutorial Languages
	err = s.Repository.AddLanguageToTutorial(
		ctx,
		req.TutorialID,
		req.Language,
		translationID,
	)
	if err != nil {
		return nil, err
	}

	// 14. Return created translation
	return translation, nil
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
	err := s.Repository.UploadVideoThumbnail(
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
	err = s.Repository.UploadVideoThumbnail(
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

	return s.Repository.GetTutorialsByRoles(ctx, cleanedRoles)
}
