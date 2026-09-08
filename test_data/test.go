package test_data

import (
	"context"
	"fmt"
	"saathi-backend/config"
	"saathi-backend/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func SeedTutorials() error {

	collection := config.DB.Collection("tutorials")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tutorials := []interface{}{

		model.Tutorial{Video_ID: "00000000-0000-0000-0000-000000000001", AppName: "saathi", Translations: map[string]model.Translation{"en": {VideoTitle: "Introduction to Redis", VideoDescription: "Learn the basics of Redis", Video_Bucket: "455719_Iceland_Iceland83_1280x720.mp4", TitleImage: "redis_image", Duration: "01:00"}}, Version: "1.0", Roles: []string{"admin", "developer"}, IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},

		model.Tutorial{Video_ID: "00000000-0000-0000-0000-000000000002", AppName: "saathi", Translations: map[string]model.Translation{"en": {VideoTitle: "Introduction to Go", VideoDescription: "Learn the basics of Golang", Video_Bucket: "455719_Iceland_Iceland83_1280x720.mp4", TitleImage: "go_image", Duration: "01:00"}}, Version: "1.0", Roles: []string{"admin", "developer"}, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},

		model.Tutorial{Video_ID: "00000000-0000-0000-0000-000000000003", AppName: "saathi", Translations: map[string]model.Translation{"en": {VideoTitle: "MongoDB Introduction", VideoDescription: "Learn MongoDB fundamentals", Video_Bucket: "455719_Iceland_Iceland83_1280x720.mp4", TitleImage: "mongodb.jpg", Duration: "01:00"}}, Version: "1.0", Roles: []string{"admin", "developer", "manager"}, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},

		model.Tutorial{Video_ID: "00000000-0000-0000-0000-000000000004", AppName: "saathi", Translations: map[string]model.Translation{"en": {VideoTitle: "Kafka Basics", VideoDescription: "Learn Apache Kafka", Video_Bucket: "455719_Iceland_Iceland83_1280x720.mp4", TitleImage: "kafka.jpg", Duration: "01:00"}}, Version: "1.0", Roles: []string{"admin", "manager"}, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	result, err := collection.InsertMany(ctx, tutorials)

	if err != nil {
		return err
	}

	fmt.Printf("Inserted %d tutorials\n", len(result.InsertedIDs))

	return nil
}
