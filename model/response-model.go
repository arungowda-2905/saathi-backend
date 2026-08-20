package model

import "time"

type TutorialRoleResponse struct {
	VideoTitle string    `bson:"video_title" json:"video_title"`
	AppName    string    `bson:"app_name" json:"app_name"`
	VideoLink  string    `bson:"video_link" json:"video_link"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
