package entity

import "time"

type Transaction struct {
TransactionID string    `gorm:"column:transaction_id;primaryKey"`
TraceID       string    `gorm:"column:trace_id;index;not null;default:''"`
AccountID     string    `gorm:"column:account_id;not null"`
MerchantID    string    `gorm:"column:merchant_id;not null"`
Amount        float64   `gorm:"column:amount;type:decimal(15,2);not null"`
Status        string    `gorm:"column:status;type:transaction_status;not null;default:PENDING"`
CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

Account  *Account  `gorm:"foreignKey:AccountID;references:AccountID"`
Merchant *Merchant `gorm:"foreignKey:MerchantID;references:MerchantID"`
}

func (Transaction) TableName() string {
return "transactions"
}
