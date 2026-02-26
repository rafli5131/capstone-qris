package entity

import "time"

type Account struct {
AccountID string    `gorm:"column:account_id;primaryKey"`
Balance   float64   `gorm:"column:balance;type:decimal(15,2);not null;default:0"`
Currency  string    `gorm:"column:currency;not null;default:IDR"`
Version   int       `gorm:"column:version;not null;default:0"`
CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Account) TableName() string {
return "accounts"
}
