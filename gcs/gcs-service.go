package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type Service struct {
	client *storage.Client
}

func NewService(ctx context.Context) (*Service, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

func (s *Service) UploadVideo(
	ctx context.Context,
	bucketName string,
	fileName string,
	file io.Reader,
) error {

	writer := s.client.
		Bucket(bucketName).
		Object(fileName).
		NewWriter(ctx)

	_, err := io.Copy(writer, file)
	if err != nil {
		writer.Close()
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
