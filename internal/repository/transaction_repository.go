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
