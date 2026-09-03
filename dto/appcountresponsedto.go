package dto

type AppCountResponseDTO struct {
	AppName string `json:"app_name" bson:"app_name"`
	Count   int64  `json:"count" bson:"count"`
}
