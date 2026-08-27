package dto

import "github.com/google/uuid"

type TutorialResponseDTO struct {
	Video_ID         uuid.UUID `bson:"video_id" json:"video_id"`
	AppName          string    `bson:"app_name" json:"app_name"`
	VideoTitle       string    `bson:"video_title" json:"video_title"`
	VideoDescription string    `bson:"video_description" json:"video_description"`
}
