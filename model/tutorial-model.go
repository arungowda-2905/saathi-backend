package model

import "time"

type Tutorial struct {
	ID      string `bson:"_id" json:"id"`
	VideoID string `bson:"video_id" json:"videoId"`
	AppName string `bson:"app_name" json:"appName"`

	Version  string   `bson:"version" json:"version"`
	Roles    []string `bson:"roles" json:"roles"`
	IsActive bool     `bson:"is_active" json:"isActive"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type TutorialTranslation struct {
	ID         string `bson:"_id" json:"id"`
	TutorialID string `bson:"tutorial_id" json:"tutorialId"`
	Language   string `bson:"language" json:"language"`

	VideoTitle       string `bson:"video_title" json:"videoTitle"`
	VideoDescription string `bson:"video_description" json:"videoDescription"`
	VideoBucket      string `bson:"video_bucket,omitempty" json:"videoBucket,omitempty"`
	TitleImage       string `bson:"title_image,omitempty" json:"titleImage,omitempty"`
	Duration         string `bson:"duration,omitempty" json:"duration,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
