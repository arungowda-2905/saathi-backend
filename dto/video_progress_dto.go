package dto

type CreateVideoProgressRequest struct {
	PositionSeconds int  `json:"position_seconds"`
	Completed       bool `json:"completed"`
}
type VideoFeedbackRequest struct {
    Rating string `json:"rating"`
}