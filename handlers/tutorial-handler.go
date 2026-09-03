package handlers

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"saathi-backend/model"
	tutorialservice "saathi-backend/tutorial-service"
)

type TutorialHandler struct {
	service *tutorialservice.TutorialService
}

func NewTutorialHandler(service *tutorialservice.TutorialService) *TutorialHandler {
	return &TutorialHandler{
		service: service,
	}
}
func (h *TutorialHandler) HandleVideoUpload(c *fiber.Ctx) error {

	// upload_id is returned by POST /v1/upload
	uploadID := c.FormValue("upload_id")

	if uploadID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Upload ID is required",
		})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	applicationName := c.FormValue("applicationName")
	assignToRole := c.FormValue("assignToRole")

	roles := parseRoles(assignToRole)

	if title == "" ||
		description == "" ||
		applicationName == "" ||
		len(roles) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All metadata fields are required",
		})
	}

	tutorial := model.Tutorial{
		AppName:          applicationName,
		VideoTitle:       title,
		VideoDescription: description,
		Roles:            roles,
		Version:          "1.0",
		IsActive:         true,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	createdTutorial, err := h.service.CreateNewTutorial(
		ctx,
		tutorial,
		uploadID,
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

func (h *TutorialHandler) GetTutorialByID(c *fiber.Ctx) error {

	vid := c.Params("videoId")

	if vid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video ID is required",
		})
	}

	// userRole := c.Get("X-Role")

	// if userRole == "" {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": "User role is required",
	// 	})
	// }

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

	userRole := c.Get("X-Role")

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

	tutorials, err := h.service.GetDetailsByRole(
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

	if len(tutorials) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No tutorials found for this role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tutorials)
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

	// Upload video + thumbnail using ONE UUID
	uploadID, err := h.service.UploadVideoAndThumbnail(
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

	// Return ONLY the UUID.
	// Do not return actual GCS file names.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Video and thumbnail uploaded successfully",
		"upload_id": uploadID,
	})
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
