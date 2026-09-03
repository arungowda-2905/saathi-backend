package gcs

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Service struct {
	client *storage.Client
}

func NewService(ctx context.Context) (*Service, error) {
	if ctx == nil {
		return nil, fmt.Errorf("GCS context is required")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

func (s *Service) Close() error {
	if s == nil || s.client == nil {
		return nil
	}

	return s.client.Close()
}

func (s *Service) UploadVideoThumbnail(
	ctx context.Context,
	bucketName string,
	fileName string,
	file io.Reader,
) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("GCS service is not initialized")
	}
	if strings.TrimSpace(bucketName) == "" {
		return fmt.Errorf("GCS bucket name is required")
	}
	if strings.TrimSpace(fileName) == "" {
		return fmt.Errorf("GCS object name is required")
	}
	if file == nil {
		return fmt.Errorf("video file is required")
	}

	writer := s.client.
		Bucket(bucketName).
		Object(fileName).
		NewWriter(ctx)
	writer.ContentType = contentType(fileName)
	writer.CacheControl = "public, max-age=3600"

	_, err := io.Copy(writer, file)
	if err != nil {
		_ = writer.Close()
		return fmt.Errorf("failed to upload video: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}

func (s *Service) GetVideo(
	ctx context.Context,
	bucketName string,
	objectName string,
) ([]byte, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("GCS service is not initialized")
	}
	if strings.TrimSpace(bucketName) == "" {
		return nil, fmt.Errorf("GCS bucket name is required")
	}
	if strings.TrimSpace(objectName) == "" {
		return nil, fmt.Errorf("GCS object name is required")
	}

	reader, err := s.client.
		Bucket(bucketName).
		Object(objectName).
		NewReader(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create GCS reader: %w",
			err,
		)
	}

	defer reader.Close()

	videoBytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to read video from GCS: %w",
			err,
		)
	}

	return videoBytes, nil
}

func (s *Service) GetUploadFile(
	ctx context.Context,
	bucketName string,
	prefix string,
	fileType string,
) (string, error) {
	if s == nil || s.client == nil {
		return "", fmt.Errorf("GCS service is not initialized")
	}
	if strings.TrimSpace(bucketName) == "" {
		return "", fmt.Errorf("GCS bucket name is required")
	}
	if strings.TrimSpace(prefix) == "" || strings.TrimSpace(fileType) == "" {
		return "", fmt.Errorf("GCS upload prefix and file type are required")
	}

	query := &storage.Query{Prefix: strings.TrimSuffix(prefix, "/")}
	it := s.client.Bucket(bucketName).Objects(ctx, query)

	for {
		object, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to list GCS upload files: %w", err)
		}
		if strings.HasPrefix(object.Name, strings.TrimSuffix(prefix, "/")) {
			return object.Name, nil
		}
	}

	return "", fmt.Errorf("%s not found for upload UUID", fileType)
}

func contentType(fileName string) string {
	if detected := mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName))); detected != "" {
		return detected
	}

	return "application/octet-stream"
}
