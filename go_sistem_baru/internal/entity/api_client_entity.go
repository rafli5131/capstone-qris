package entity

import "time"

type ApiClient struct {
	ClientID     string    `gorm:"column:client_id;primaryKey"`
	ClientSecret string    `gorm:"column:client_secret;not null"`
	Status       string    `gorm:"column:status;type:api_client_status;not null;default:ACTIVE"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ApiClient) TableName() string {
	return "api_clients"
}
