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

	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video file is required",
		})
	}

	if err := validateVideoFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded video",
		})
	}
	defer src.Close()

	if err := validateVideoContent(src); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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

		Roles: roles,

		Version:    "1.0",
		TitleImage: "",
		IsActive:   true,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	videoFileName, err := h.service.UploadVideo(
		c.Context(),
		file.Filename,
		src,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload video",
		})
	}

	createdTutorial, err := h.service.CreateNewTutorial(
		ctx,
		tutorial,
		videoFileName,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tutorial",
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
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Video file is required",
		})
	}

	if err := validateVideoFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded video",
		})
	}
	defer src.Close()

	if err := validateVideoContent(src); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fileName, err := h.service.UploadVideo(
		c.Context(),
		file.Filename,
		src,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload video",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Video uploaded successfully",
		"fileName": fileName,
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
