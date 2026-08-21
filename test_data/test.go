package test_data

import (
	"context"
	"fmt"
	"time"

	"saathi-backend/config"
	"saathi-backend/model"
)

func SeedTutorials() error {

	collection := config.DB.Collection("tutorials")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tutorials := []interface{}{

		model.Tutorial{

			Video_Bucket:     "455719_Iceland_Iceland83_1280x720",
			AppName:          "saathi",
			VideoTitle:       "Introduction to Redis",
			VideoDescription: "Learn the basics of Redis",
			Version:          "1.0",
			TitleImage:       "redis_image",
			Roles:            []string{"admin", "developer"},
			IsActive:         false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		model.Tutorial{

			Video_Bucket:     "455719_Iceland_Iceland83_1280x720",
			AppName:          "saathi",
			VideoTitle:       "Introduction to Go",
			VideoDescription: "Learn the basics of Golang",
			Version:          "1.0",
			TitleImage:       "go_image",
			Roles:            []string{"admin", "developer"},
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		model.Tutorial{

			Video_Bucket:     "455719_Iceland_Iceland83_1280x720",
			AppName:          "saathi",
			VideoTitle:       "MongoDB Introduction",
			VideoDescription: "Learn MongoDB fundamentals",
			Version:          "1.0",
			TitleImage:       "https://storage.googleapis.com/mybucket/mongodb.jpg",
			Roles:            []string{"admin", "developer", "manager"},
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		model.Tutorial{

			Video_Bucket:     "455719_Iceland_Iceland83_1280x720",
			AppName:          "saathi",
			VideoTitle:       "Kafka Basics",
			VideoDescription: "Learn Apache Kafka",
			Version:          "1.0",
			TitleImage:       "https://storage.googleapis.com/mybucket/kafka.jpg",
			Roles:            []string{"admin", "manager"},
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}

	result, err := collection.InsertMany(ctx, tutorials)

	if err != nil {
		return err
	}

	fmt.Printf("Inserted %d tutorials\n", len(result.InsertedIDs))

	return nil
}
