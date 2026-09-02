package model

import "time"

type VideoProgress struct {
	ID              string    `json:"id" bson:"_id"`
	UserID          string    `json:"user_id" bson:"user_id"`
	VideoID         string    `json:"video_id" bson:"video_id"`
	PositionSeconds int       `json:"position_seconds" bson:"position_seconds"`
	Completed       bool      `json:"completed" bson:"completed"`
	LastWatchedAt   time.Time `json:"last_watched_at" bson:"last_watched_at"`
	Rating          string    `json:"rating,omitempty" bson:"rating,omitempty"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}