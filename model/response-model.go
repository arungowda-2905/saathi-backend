package model

import "go.mongodb.org/mongo-driver/v2/bson"

type TutorialDetail struct {
	ID               bson.ObjectID `bson:"_id" json:"id"`
	AppName          string        `bson:"app_name" json:"app_name"`
	VideoTitle       string        `bson:"video_title" json:"video_title"`
	VideoDescription string        `bson:"video_description" json:"video_description"`
}
