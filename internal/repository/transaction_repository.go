package repository

import (
	"capstone-qris/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(db *gorm.DB, log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Repository: Repository[entity.Transaction]{DB: db},
		Log:        log,
	}
}

func (r *TransactionRepository) FindByTransactionID(db *gorm.DB, transactionID string) (*entity.Transaction, error) {
	var tx entity.Transaction
	if err := db.Where("transaction_id = ?", transactionID).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) UpdateStatus(db *gorm.DB, transactionID, status string) error {
	return db.Model(&entity.Transaction{}).
		Where("transaction_id = ?", transactionID).
		Update("status", status).Error
}

func (r *TransactionRepository) FindTransactionsByMerchantID(db *gorm.DB, merchantID string) ([]entity.Transaction, error) {
	var txs []entity.Transaction
	if err := db.Where("merchant_id = ?", merchantID).
		Preload("Account").Preload("Merchant").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *TransactionRepository) SumIncomeByMerchantID(db *gorm.DB, merchantID string) (float64, error) {
	var total float64
	if err := db.Model(&entity.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("merchant_id = ? AND status = ?", merchantID, "SUCCESS").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *TransactionRepository) FindAll(db *gorm.DB) ([]entity.Transaction, error) {
	var txs []entity.Transaction
	if err := db.Preload("Account").Preload("Merchant").Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *TransactionRepository) Update(db *gorm.DB, tx *entity.Transaction) error {
	return db.Save(tx).Error
}
