package dto

type CreateTutorialRequest struct {
	AppName             string   `json:"appName"`
	TutorialTitle       string   `json:"tutorialTitle"`
	TutorialDescription string   `json:"tutorialDescription"`
	VideoBucket         string   `json:"videoBucket"`
	ThumbnailImage      string   `json:"titleImage"`
	Duration            string   `json:"duration"`
	Version             string   `json:"version"`
	Roles               []string `json:"roles"`
	IsActive            bool     `json:"isActive"`
}

type CreateTranslationRequest struct {
	TutorialID          string `json:"tutorialId"`
	TutorialTitle       string `json:"tutorialTitle"`
	TutorialDescription string `json:"tutorialDescription"`
	VideoBucket         string `json:"videoBucket"`
	ThumbnailImage      string `json:"titleImage"`
	Duration            string `json:"duration"`
	Language            string `json:"language"`
}

type AppCountResponseDTO struct {
	AppName string `json:"app_name" bson:"app_name"`
	Count   int64  `json:"count" bson:"count"`
}
