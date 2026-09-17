package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Translation struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TranslationID       string             `bson:"translation_id" json:"translationId"`
	TutorialTitle       string             `bson:"tutorial_title" json:"tutorialTitle"`
	TutorialDescription string             `bson:"tutorial_description" json:"tutorialDescription"`
	VideoBucket         string             `bson:"video_bucket,omitempty" json:"videoBucket,omitempty"`
	ThumbnailImage      string             `bson:"title_image" json:"titleImage"`
	Duration            string             `bson:"duration" json:"duration"`
	Language            string             `bson:"language" json:"language"`
}

type Tutorial struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TutorialID          string             `bson:"tutorial_id" json:"tutorialId"`
	AppName             string             `bson:"app_name" json:"appName"`
	TutorialTitle       string             `bson:"tutorial_title" json:"tutorialTitle"`
	TutorialDescription string             `bson:"tutorial_description" json:"tutorialDescription"`
	VideoBucket         string             `bson:"video_bucket,omitempty" json:"videoBucket,omitempty"`
	ThumbnailImage      string             `bson:"title_image" json:"titleImage"`
	Duration            string             `bson:"duration" json:"duration"`

	Languages map[string]string `bson:"languages" json:"languages"`

	Version  string   `bson:"version" json:"version"`
	Roles    []string `bson:"roles" json:"roles"`
	IsActive bool     `bson:"is_active" json:"isActive"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
