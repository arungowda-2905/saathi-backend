package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Tutorial struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AppName          string        `bson:"app_name" json:"app_name"`
	VideoTitle       string        `bson:"video_title" json:"video_title"`
	VideoDescription string        `bson:"video_description" json:"video_description"`
	VideoLink        string        `bson:"video_link" json:"video_link"`
	Version          string        `bson:"version" json:"version"`
	TitleImage       string        `bson:"title_image" json:"title_image"`
	Roles            []string      `bson:"roles" json:"roles"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
}
