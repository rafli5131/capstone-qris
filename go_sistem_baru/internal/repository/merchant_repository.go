package repository

import (
	"capstone-qris/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MerchantRepository struct {
	Repository[entity.Merchant]
	Log *logrus.Logger
}

func NewMerchantRepository(db *gorm.DB, log *logrus.Logger) *MerchantRepository {
	return &MerchantRepository{
		Repository: Repository[entity.Merchant]{DB: db},
		Log:        log,
	}
}

func (r *MerchantRepository) FindByMerchantID(db *gorm.DB, merchantID string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	if err := db.Where("merchant_id = ? AND is_active = true", merchantID).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *MerchantRepository) FindByMerchantIDIncludingInactive(db *gorm.DB, merchantID string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	if err := db.Where("merchant_id = ?", merchantID).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}
