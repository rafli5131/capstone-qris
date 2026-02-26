package model

// PaymentRequest is the body for POST /api/qris/payment
type PaymentRequest struct {
	InquiryID     string  `json:"inquiry_id"      validate:"required"`
	UserID        string  `json:"user_id"         validate:"required"`
	Amount        float64 `json:"amount"          validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method"  validate:"required,oneof=balance"`
	Pincode       string  `json:"pincode"         validate:"required,len=6"`
}

// PaymentResponse is the immediate response for POST /api/qris/payment
type PaymentResponse struct {
	Status              string `json:"status"`
	TransactionID       string `json:"transaction_id"`
	Message             string `json:"message"`
	EstimatedCompletion string `json:"estimated_completion"`
}
