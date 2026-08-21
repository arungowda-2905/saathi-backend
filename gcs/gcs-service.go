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
