package model

import "time"

type ApiClientCreateRequest struct {
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Status       string `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}

type ApiClientUpdateRequest struct {
	ClientSecret string `json:"client_secret" validate:"omitempty"`
	Status       string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

type ApiClientResponse struct {
	ClientID  string    `json:"client_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionUpdateRequest struct {
	Amount float64 `json:"amount" validate:"omitempty,gt=0"`
	Status string  `json:"status" validate:"omitempty,oneof=PENDING SUCCESS FAILED"`
}

type TransactionListItem struct {
	TransactionID string    `json:"transaction_id"`
	TraceID       string    `json:"trace_id"`
	AccountID     string    `json:"account_id"`
	MerchantID    string    `json:"merchant_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
