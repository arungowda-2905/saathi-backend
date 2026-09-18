package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"saathi-backend/dto"
	"saathi-backend/gcs"
	tutorialservice "saathi-backend/tutorial-service"

	"github.com/abema/go-mp4"
	"github.com/gofiber/fiber/v2"
)

type TutorialHandler struct {
	Service *tutorialservice.TutorialService
}

type TutorialRepository struct {
	gcsService *gcs.Service
}

func NewTutorialRepository(gcsService *gcs.Service) *TutorialRepository {
	return &TutorialRepository{
		gcsService: gcsService,
	}
}

func NewTutorialHandler(
	service *tutorialservice.TutorialService,
) *TutorialHandler {
	return &TutorialHandler{
		Service: service,
	}
}

func (h *TutorialHandler) CreateTutorial(c *fiber.Ctx) error {

	var req dto.CreateTutorialRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(
		c.Context(),
		10*time.Second,
	)
	defer cancel()

	tutorial, err := h.Service.CreateTutorial(ctx, req)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success":     true,
		"message":     "Tutorial created successfully",
		"tutorial_id": tutorial,
	})
}

func (h *TutorialHandler) CreateTranslation(
	c *fiber.Ctx,
) error {

	// 1. Parse request body

	var req dto.CreateTranslationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// 2. Create context with timeout

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// 3. Call service

	_, err := h.Service.CreateTranslation(
		ctx,
		req,
	)

	if err != nil {

		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{

			"success": false,
			"message": err.Error(),
		})
	}

	// -----------------------------------------
	// 4. Return response
	// -----------------------------------------

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{

		"success": true,
		"message": "Translation created successfully",
		// "data":    translation,
	})
}

func (h *TutorialHandler) UploadVideo(c *fiber.Ctx) error {

	// Get video file
	videoFile, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video file is required",
		})
	}

	// Get thumbnail file
	thumbnailFile, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thumbnail image is required",
		})
	}

	// Validate video
	if err := validateVideoFile(videoFile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate thumbnail
	if err := validateThumbnailFile(thumbnailFile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Open video
	videoSrc, err := videoFile.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded video",
		})
	}
	defer videoSrc.Close()

	if err := validateVideoContent(videoSrc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	duration, err := videoDuration(videoSrc)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	thumbnailSrc, err := thumbnailFile.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded thumbnail",
		})
	}
	defer thumbnailSrc.Close()

	if err := validateThumbnailContent(thumbnailSrc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	videoBucketUUID, thumbnailUUID, err := h.Service.UploadVideoAndThumbnail(
		c.Context(),
		videoFile.Filename,
		videoSrc,
		thumbnailFile.Filename,
		thumbnailSrc,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload video and thumbnail",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":           "Video and thumbnail uploaded successfully",
		"video_bucket_uuid": videoBucketUUID,
		"thumbnail_uuid":    thumbnailUUID,
		"duration":          duration,
	})
}

func validateVideoFile(file *multipart.FileHeader) error {
	if file.Size == 0 {
		return errors.New("Video file cannot be empty")
	}

	switch strings.ToLower(filepath.Ext(file.Filename)) {
	case ".mp4", ".mov", ".webm", ".avi", ".mkv":
		return nil
	default:
		return errors.New("Unsupported video format")
	}
}

func validateThumbnailFile(file *multipart.FileHeader) error {
	if file.Size == 0 {
		return errors.New("Thumbnail image cannot be empty")
	}

	switch strings.ToLower(filepath.Ext(file.Filename)) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return nil
	default:
		return errors.New("Unsupported thumbnail format")
	}
}

func validateVideoContent(file multipart.File) error {
	const sniffSize = 512

	header := make([]byte, sniffSize)
	read, err := file.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.New("Unable to read video file")
	}

	contentType := http.DetectContentType(header[:read])
	if !strings.HasPrefix(contentType, "video/") {
		return errors.New("Uploaded file is not a valid video")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return errors.New("Unable to process video file")
	}

	return nil
}

func videoDuration(file multipart.File) (string, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("unable to seek video: %w", err)
	}

	var movieHeader *mp4.Mvhd

	_, err := mp4.ReadBoxStructure(file, func(handle *mp4.ReadHandle) (interface{}, error) {
		fmt.Printf(
			"BOX: type=%v size=%d\n",
			handle.BoxInfo.Type,
			handle.BoxInfo.Size,
		)

		if handle.BoxInfo.Type == mp4.BoxTypeMvhd() {
			box, _, err := handle.ReadPayload()
			if err != nil {
				return nil, fmt.Errorf("failed to read mvhd: %w", err)
			}

			var ok bool
			movieHeader, ok = box.(*mp4.Mvhd)
			if !ok {
				return nil, errors.New("invalid mvhd box")
			}

			return movieHeader, nil
		}

		if handle.BoxInfo.Type == mp4.BoxTypeMoov() {
			_, err := handle.Expand()
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse MP4: %w", err)
	}

	if movieHeader == nil {
		return "", errors.New("mvhd box not found")
	}

	if movieHeader.Timescale == 0 {
		return "", errors.New("invalid video timescale")
	}

	duration := movieHeader.GetDuration()

	totalSeconds := int(
		(duration + uint64(movieHeader.Timescale/2)) /
			uint64(movieHeader.Timescale),
	)

	return fmt.Sprintf(
		"%02d:%02d",
		totalSeconds/60,
		totalSeconds%60,
	), nil
}

func validateThumbnailContent(file multipart.File) error {
	const sniffSize = 512

	header := make([]byte, sniffSize)

	read, err := file.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.New("Unable to read thumbnail image")
	}

	contentType := http.DetectContentType(header[:read])
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("Uploaded file is not a valid image")
	}

	// Reset reader position so it can be uploaded to GCS.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return errors.New("Unable to process thumbnail image")
	}

	return nil
}

func (h *TutorialHandler) GetDetailsByRole(c *fiber.Ctx) error {

	userRole := c.Get("X-Role")
	if userRole == "" {
		userRole = c.Query("X-Role")
	}
	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role not found",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	appCounts, err := h.Service.GetDetailsByRole(
		ctx,
		userRole,
	)

	if err != nil {

		if errors.Is(err, tutorialservice.ErrRoleNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorials",
		})
	}

	if len(appCounts) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No tutorials found for this role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(appCounts)
}
