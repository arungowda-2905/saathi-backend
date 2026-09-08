package dto

type TutorialResponseDTO struct {
	Video_ID         string `json:"video_id"`
	AppName          string `json:"app_name"`
	VideoTitle       string `json:"video_title"`
	VideoDescription string `json:"video_description"`
	Video_Bucket     string `json:"video_bucket,omitempty"`
	TitleImage       string `json:"title_image"`
	Duration         string `json:"duration"`
}
