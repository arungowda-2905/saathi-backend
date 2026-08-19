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
			AppName:          "saathi",
			VideoTitle:       "Introduction to Go",
			VideoDescription: "Learn the basics of Golang",
			VideoLink:        "https://storage.googleapis.com/mybucket/go-introduction.mp4",
			Version:          "1.0",
			TitleImage:       "https://storage.googleapis.com/mybucket/go.jpg",
			Roles:            []string{"admin", "developer"},
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		model.Tutorial{
			AppName:          "saathi",
			VideoTitle:       "MongoDB Introduction",
			VideoDescription: "Learn MongoDB fundamentals",
			VideoLink:        "https://storage.googleapis.com/mybucket/mongodb-introduction.mp4",
			Version:          "1.0",
			TitleImage:       "https://storage.googleapis.com/mybucket/mongodb.jpg",
			Roles:            []string{"admin", "developer", "manager"},
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},

		model.Tutorial{
			AppName:          "saathi",
			VideoTitle:       "Kafka Basics",
			VideoDescription: "Learn Apache Kafka",
			VideoLink:        "https://storage.googleapis.com/mybucket/kafka.mp4",
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
