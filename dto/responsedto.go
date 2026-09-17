package dto

type HandleVideoUploadRequest struct {
	AppName  string   `json:"appName"`
	Version  string   `json:"version"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"isActive"`

	Language         string `json:"language"`
	VideoTitle       string `json:"videoTitle"`
	VideoDescription string `json:"videoDescription"`
	VideoBucket      string `json:"videoBucket"`
	TitleImage       string `json:"titleImage"`
	Duration         string `json:"duration"`
}
