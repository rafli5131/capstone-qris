package repository

import (
	"capstone-qris/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountRepository struct {
	Repository[entity.Account]
	Log *logrus.Logger
}

func NewAccountRepository(db *gorm.DB, log *logrus.Logger) *AccountRepository {
	return &AccountRepository{
		Repository: Repository[entity.Account]{DB: db},
		Log:        log,
	}
}

func (r *AccountRepository) FindByAccountID(db *gorm.DB, accountID string) (*entity.Account, error) {
	var account entity.Account
	if err := db.Where("account_id = ?", accountID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByUsername(db *gorm.DB, username string) (*entity.Account, error) {
	var account entity.Account
	if err := db.Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) UpdateWithOptimisticLock(db *gorm.DB, accountID string, amount float64, currentVersion int) (int64, error) {
	result := db.Model(&entity.Account{}).
		Where("account_id = ? AND version = ?", accountID, currentVersion).
		Updates(map[string]interface{}{
			"balance": gorm.Expr("balance - ?", amount),
			"version": gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}
