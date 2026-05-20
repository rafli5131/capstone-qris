package model

import "time"

type TransactionStatusResponse struct {
TransactionID string    `json:"transaction_id"`
Status        string    `json:"status"`
FinalBalance  float64   `json:"final_balance"`
Timestamp     time.Time `json:"timestamp"`
}
