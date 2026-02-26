package model

// WebResponse is the generic top-level envelope for all API responses.
type WebResponse[T any] struct {
	Status   string    `json:"status"`
	Data     T         `json:"data,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Errors   string    `json:"errors,omitempty"`
}

// Metadata contains performance and source information.
type Metadata struct {
	LatencyMS int64  `json:"latency_ms"`
	Source    string `json:"source,omitempty"` // "cache" or "database"
}

// PaymentResponseEnvelope is a Swagger-friendly wrapper for payment responses.
type PaymentResponseEnvelope struct {
	Status string           `json:"status"`
	Data   *PaymentResponse `json:"data,omitempty"`
}

// InquiryResponseEnvelope is a Swagger-friendly wrapper for inquiry responses.
type InquiryResponseEnvelope struct {
	Status   string           `json:"status"`
	Data     *InquiryResponse `json:"data,omitempty"`
	Metadata *Metadata        `json:"metadata,omitempty"`
}

// ErrorResponse is a Swagger-friendly error payload.
type ErrorResponse struct {
	Status string `json:"status"`
	Errors string `json:"errors"`
}
