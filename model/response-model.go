package model

import "github.com/google/uuid"

type TutorialDetail struct {
	Video_ID         uuid.UUID
	AppName          string `bson:"app_name" json:"app_name"`
	VideoTitle       string `bson:"video_title" json:"video_title"`
	VideoDescription string `bson:"video_description" json:"video_description"`
}
