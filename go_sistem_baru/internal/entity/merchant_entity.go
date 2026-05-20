package entity

import "time"

type Merchant struct {
MerchantID   string    `gorm:"column:merchant_id;primaryKey"`
MerchantName string    `gorm:"column:merchant_name;not null"`
MCC          string    `gorm:"column:mcc;not null;default:''"`
City         string    `gorm:"column:city;not null;default:''"`
IsActive     bool      `gorm:"column:is_active;not null;default:true"`
CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Merchant) TableName() string {
return "merchants"
}
