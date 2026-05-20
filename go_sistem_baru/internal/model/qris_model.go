package model

// InquiryResponse is the data payload for GET /api/qris/inquiry/:qris_payload
type InquiryResponse struct {
	MerchantID   string  `json:"merchant_id"`
	MerchantName string  `json:"merchant_name"`
	TerminalID   string  `json:"terminal_id"`
	City         string  `json:"city"`
	FixedAmount  float64 `json:"fixed_amount"`
	InquiryID    string  `json:"inquiry_id"`
}
