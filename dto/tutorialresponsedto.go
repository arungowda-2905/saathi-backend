package dto

import (
	"github.com/google/uuid"
)

type TutorialResponseDTO struct {
	Video_ID         uuid.UUID `bson:"video_id" json:"video_id"`
	VideoTitle       string    `bson:"video_title" json:"video_title"`
	VideoDescription string    `bson:"video_description" json:"video_description"`
	TitleImage       string    `json:"title_image" bson:"title_image"`
	Duration         string    `json:"duration" bson:"duration"`
	CreatedAt        string    `bson:"created_at" json:"created_at"`
}
