package model

// MerchantResponse is the data payload for merchant creation from QR code image.
type MerchantResponse struct {
	MerchantID   string `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	City         string `json:"city"`
	MCC          string `json:"mcc"`
	IsActive     bool   `json:"is_active"`
	IsNew        bool   `json:"is_new"`
}

// MerchantResponseEnvelope is a Swagger-friendly wrapper for merchant responses.
type MerchantResponseEnvelope struct {
	Status string            `json:"status"`
	Data   *MerchantResponse `json:"data,omitempty"`
}
