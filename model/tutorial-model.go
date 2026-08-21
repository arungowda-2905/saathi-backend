package model

import (
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Tutorial struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Video_ID         uuid.UUID     `bson:"video_id" json:"video_id" validate:"required"`
	Video_Bucket     string        `bson:"video_bucket" json:"video_bucket" validate:"required"`
	AppName          string        `bson:"app_name" json:"app_name" validate:"required"`
	VideoTitle       string        `bson:"video_title" json:"video_title" validate:"required"`
	VideoDescription string        `bson:"video_description" json:"video_description" validate:"required"`
	Version          string        `bson:"version" json:"version" validate:"required"`
	TitleImage       string        `bson:"title_image" json:"title_image" validate:"required"`
	Roles            []string      `bson:"roles" json:"roles" validate:"required,min=1"`
	IsActive         bool          `bson:"is_active" json:"is_active"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
}
