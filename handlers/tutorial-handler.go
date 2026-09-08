package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mp4 "github.com/abema/go-mp4"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

type TutorialHandler struct {
	service *tutorialservice.TutorialService
}

type AddTranslationRequest struct {
	Language         string `json:"language"`
	VideoTitle       string `json:"video_title"`
	VideoDescription string `json:"video_description"`
	VideoBucket      string `json:"video_bucket"`
	TitleImage       string `json:"title_image"`
	Duration         string `json:"duration"`
}

func NewTutorialHandler(service *tutorialservice.TutorialService) *TutorialHandler {
	return &TutorialHandler{
		service: service,
	}
}
func (h *TutorialHandler) HandleVideoUpload(c *fiber.Ctx) error {

	videoBucketUUID := c.FormValue("video_bucket_uuid")
	thumbnailUUID := c.FormValue("thumbnail_uuid")

	if videoBucketUUID == "" || thumbnailUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video bucket UUID and thumbnail UUID are required",
		})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	applicationName := c.FormValue("applicationName")
	assignToRole := c.FormValue("assignToRole")
	language := c.FormValue("language")
	duration := strings.TrimSpace(c.FormValue("duration"))
	roles := parseRoles(assignToRole)

	if title == "" ||
		description == "" ||
		applicationName == "" ||
		language == "" ||
		duration == "" ||
		len(roles) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All metadata fields are required",
		})
	}
	if !isValidDuration(duration) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Duration must use MM:SS format",
		})
	}
	language = strings.ToLower(language)
	if !isSupportedLanguage(language) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid language",
		})
	}

	translation := model.Translation{
		VideoTitle:       title,
		VideoDescription: description,
		Video_Bucket:     videoBucketUUID,
		TitleImage:       thumbnailUUID,
		Duration:         duration,
	}
	tutorial := model.Tutorial{
		AppName:      applicationName,
		Translations: map[string]model.Translation{language: translation},
		Roles:        roles,
		Version:      "1.0",
		IsActive:     true,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	createdTutorial, err := h.service.CreateNewTutorial(
		ctx,
		tutorial,
		videoBucketUUID,
		thumbnailUUID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Video uploaded successfully",
		"video_id": createdTutorial.Video_ID,
	})
}

func (h *TutorialHandler) AddTutorialTranslation(c *fiber.Ctx) error {
	videoID := strings.TrimSpace(c.Params("videoId"))
	if videoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Video ID is required"})
	}

	var request AddTranslationRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	request.Language = strings.ToLower(strings.TrimSpace(request.Language))

	if !isSupportedLanguage(request.Language) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid language"})
	}
	if request.VideoTitle == "" || request.VideoDescription == "" || request.VideoBucket == "" || request.TitleImage == "" || !isValidDuration(request.Duration) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Translation fields are invalid"})
	}

	translation := model.Translation{
		VideoTitle:       request.VideoTitle,
		VideoDescription: request.VideoDescription,
		Video_Bucket:     request.VideoBucket,
		TitleImage:       request.TitleImage,
		Duration:         request.Duration,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.service.AddTranslation(ctx, videoID, request.Language, translation); err != nil {
		if errors.Is(err, tutorialservice.ErrInvalidVideoID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid video ID"})
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tutorial not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add translation"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Translation added successfully",
		"video_id": videoID,
		"language": request.Language,
	})
}

func (h *TutorialHandler) GetTutorialByID(c *fiber.Ctx) error {

	vid := c.Params("videoId")

	if vid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	videoBytes, err := h.service.GetTutorialByID(
		vid,
		ctx,
	)
	// userRole,

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tutorial not found",
			})
		}

		if errors.Is(err, tutorialservice.ErrInvalidVideoID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid video ID",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorial",
		})
	}

	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Length", strconv.Itoa(len(videoBytes)))

	return c.Send(videoBytes)
}

func (h *TutorialHandler) GetDetailsByRole(c *fiber.Ctx) error {

	//userRole := c.Get("X-Role")
	userRole := "admin" // Hardcoded for testing purposes. Remove this line in production.
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

	appCounts, err := h.service.GetDetailsByRole(
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

	// Validate actual video content
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

	// Open thumbnail
	thumbnailSrc, err := thumbnailFile.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded thumbnail",
		})
	}
	defer thumbnailSrc.Close()

	// Validate actual thumbnail content
	if err := validateThumbnailContent(thumbnailSrc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Upload video and thumbnail using independent UUID prefixes.
	videoBucketUUID, thumbnailUUID, err := h.service.UploadVideoAndThumbnail(
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

func videoDuration(file multipart.File) (string, error) {
	var movieHeader *mp4.Mvhd

	_, err := mp4.ReadBoxStructure(file, func(handle *mp4.ReadHandle) (interface{}, error) {
		if handle.BoxInfo.Type == mp4.BoxTypeMvhd() {
			box, _, err := handle.ReadPayload()
			if err != nil {
				return nil, err
			}

			var ok bool
			movieHeader, ok = box.(*mp4.Mvhd)
			if !ok {
				return nil, errors.New("invalid video movie header")
			}

			return movieHeader, nil
		}

		_, err := handle.Expand()
		return nil, err
	})
	if err != nil {
		return "", errors.New("Unable to read video duration")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("Unable to process video file")
	}
	fmt.Println("movieHeader:", movieHeader)
	if movieHeader == nil || movieHeader.Timescale == 0 {
		return "", errors.New("Unable to read video duration")
	}

	totalSeconds := int((movieHeader.GetDuration() + uint64(movieHeader.Timescale/2)) / uint64(movieHeader.Timescale))
	return fmt.Sprintf("%02d:%02d", totalSeconds/60, totalSeconds%60), nil
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

func (h *TutorialHandler) GetDetailsByAppName(c *fiber.Ctx) error {

	appName := c.Query("appName")

	if strings.TrimSpace(appName) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "App name is required",
		})
	}

	language := c.Get("X-Language")

	if strings.TrimSpace(language) == "" {
		language = "en"
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	tutorials, err := h.service.GetDetailsByAppName(
		ctx,
		appName,
		language,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorials",
		})
	}

	if len(tutorials) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No tutorials found for this app",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tutorials)
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

func parseRoles(value string) []string {
	seen := make(map[string]struct{})
	roles := make([]string, 0)

	for _, role := range strings.Split(value, ",") {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		if _, exists := seen[role]; exists {
			continue
		}

		seen[role] = struct{}{}
		roles = append(roles, role)
	}

	return roles
}

func isSupportedLanguage(language string) bool {
	switch language {
	case "en", "hi", "kn", "es", "fr", "de", "zh", "ja", "ar", "pt", "ru":
		return true
	default:
		return false
	}
}

func isValidDuration(value string) bool {
	parts := strings.Split(value, ":")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) != 2 {
		return false
	}

	minutes, minuteErr := strconv.Atoi(parts[0])
	seconds, secondErr := strconv.Atoi(parts[1])
	return minuteErr == nil && secondErr == nil && minutes >= 0 && seconds >= 0 && seconds < 60
}

func (h *TutorialHandler) GetTutorialThumbnailByID(c *fiber.Ctx) error {

	vid := c.Params("videoId")

	if vid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	imageBytes, contentType, err := h.service.GetTutorialThumbnailByID(
		vid,
		ctx,
	)

	if err != nil {

		fmt.Println("ERROR GetTutorialThumbnailByID:", err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tutorial not found",
			})
		}

		if errors.Is(err, tutorialservice.ErrInvalidVideoID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid video ID",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tutorial thumbnail",
		})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Length", strconv.Itoa(len(imageBytes)))
	c.Set("Cache-Control", "public, max-age=3600")

	return c.Send(imageBytes)
}
