package model

type MerchantIncomeResponse struct {
	MerchantID       string                `json:"merchant_id"`
	MerchantName     string                `json:"merchant_name"`
	TotalRevenue     float64               `json:"total_revenue"`
	TransactionCount int                   `json:"transaction_count"`
	Transactions     []TransactionListItem `json:"transactions"`
}
