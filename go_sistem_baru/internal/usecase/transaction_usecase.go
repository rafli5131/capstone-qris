package usecase

import (
	"context"
	"errors"

	"capstone-qris/internal/model"
	"capstone-qris/internal/model/converter"
	"capstone-qris/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TransactionUseCase handles transaction status queries.
type TransactionUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	TransactionRepo *repository.TransactionRepository
	AccountRepo     *repository.AccountRepository
}

func NewTransactionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		DB:              db,
		Log:             log,
		TransactionRepo: transactionRepo,
		AccountRepo:     accountRepo,
	}
}

// GetStatus returns the current status of a transaction along with the account's final balance.
func (uc *TransactionUseCase) GetStatus(ctx context.Context, transactionID string) (*model.TransactionStatusResponse, error) {
	tx, err := uc.TransactionRepo.FindByTransactionID(uc.DB.WithContext(ctx), transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "transaction not found")
		}
		uc.Log.WithError(err).Warn("failed to find transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToTransactionStatusResponse(tx), nil
}
