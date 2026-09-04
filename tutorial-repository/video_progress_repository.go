package tutorialrepository

import (
	"context"

	"saathi-backend/config"
	"saathi-backend/model"
"time"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var videoProgressCollection *mongo.Collection

func InitVideoProgressRepository() {
	videoProgressCollection = config.DB.Collection("video_progress")
}

func CreateVideoProgress(progress *model.VideoProgress) error {
	_, err := videoProgressCollection.InsertOne(
		context.Background(),
		progress,
	)

	return err
}

func GetVideoProgress(userID, videoID string) (*model.VideoProgress, error) {
	var progress model.VideoProgress

	err := videoProgressCollection.FindOne(
		context.Background(),
		bson.M{
			"user_id":  userID,
			"video_id": videoID,
		},
	).Decode(&progress)

	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func CreateVideoProgressIndex() error {
	_, err := videoProgressCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "video_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	return err
}
func UpdateVideoFeedback(userID, videoID, rating string) error {
    _, err := videoProgressCollection.UpdateOne(
        context.Background(),
        bson.M{
            "user_id":  userID,
            "video_id": videoID,
        },
        bson.M{
            "$set": bson.M{
                "rating":     rating,
                "updated_at": time.Now(),
            },
        },
    )

    return err
}

func UpdateVideoProgress(
	userID string,
	videoID string,
	positionSeconds int,
	completed bool,
) error {

	_, err := videoProgressCollection.UpdateOne(
		context.Background(),
		bson.M{
			"user_id":  userID,
			"video_id": videoID,
		},
		bson.M{
			"$set": bson.M{
				"position_seconds": positionSeconds,
				"completed":        completed,
				"last_watched_at":  time.Now(),
				"updated_at":       time.Now(),
			},
		},
	)

	return err
}