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
			VideoTitle:       "Introduction to Redis",
			VideoDescription: "Learn the basics of Redis",
			VideoLink:        "https://storage.googleapis.com/mybucket/redis-introduction.mp4",
			Version:          "1.0",
			TitleImage:       "https://storage.googleapis.com/mybucket/redis.jpg",
			Roles:            []string{"admin", "developer"},
			IsActive:         false,
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
