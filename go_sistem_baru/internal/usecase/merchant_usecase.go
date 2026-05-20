package usecase

import (
	"context"
	"errors"

	"capstone-qris/internal/model"
	"capstone-qris/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MerchantUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	MerchantRepo    *repository.MerchantRepository
	TransactionRepo *repository.TransactionRepository
}

func NewMerchantUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	merchantRepo *repository.MerchantRepository,
	transactionRepo *repository.TransactionRepository,
) *MerchantUseCase {
	return &MerchantUseCase{
		DB:              db,
		Log:             log,
		MerchantRepo:    merchantRepo,
		TransactionRepo: transactionRepo,
	}
}

func (uc *MerchantUseCase) GetIncome(ctx context.Context, merchantID string) (*model.MerchantIncomeResponse, error) {
	merchant, err := uc.MerchantRepo.FindByMerchantID(uc.DB.WithContext(ctx), merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "merchant not found")
		}
		uc.Log.WithError(err).Warn("failed to find merchant")
		return nil, fiber.ErrInternalServerError
	}

	txs, err := uc.TransactionRepo.FindTransactionsByMerchantID(uc.DB.WithContext(ctx), merchantID)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to query merchant transactions")
		return nil, fiber.ErrInternalServerError
	}

	total, err := uc.TransactionRepo.SumIncomeByMerchantID(uc.DB.WithContext(ctx), merchantID)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to sum merchant income")
		return nil, fiber.ErrInternalServerError
	}

	items := make([]model.TransactionListItem, 0, len(txs))
	for _, tx := range txs {
		items = append(items, model.TransactionListItem{
			TransactionID: tx.TransactionID,
			TraceID:       tx.TraceID,
			AccountID:     tx.AccountID,
			MerchantID:    tx.MerchantID,
			Amount:        tx.Amount,
			Status:        tx.Status,
			CreatedAt:     tx.CreatedAt,
			UpdatedAt:     tx.UpdatedAt,
		})
	}

	return &model.MerchantIncomeResponse{
		MerchantID:       merchant.MerchantID,
		MerchantName:     merchant.MerchantName,
		TotalRevenue:     total,
		TransactionCount: len(items),
		Transactions:     items,
	}, nil
}

func (uc *MerchantUseCase) GetTransactions(ctx context.Context, merchantID string) ([]model.TransactionListItem, error) {
	_, err := uc.MerchantRepo.FindByMerchantID(uc.DB.WithContext(ctx), merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "merchant not found")
		}
		uc.Log.WithError(err).Warn("failed to find merchant")
		return nil, fiber.ErrInternalServerError
	}

	txs, err := uc.TransactionRepo.FindTransactionsByMerchantID(uc.DB.WithContext(ctx), merchantID)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to query merchant transactions")
		return nil, fiber.ErrInternalServerError
	}

	items := make([]model.TransactionListItem, 0, len(txs))
	for _, tx := range txs {
		items = append(items, model.TransactionListItem{
			TransactionID: tx.TransactionID,
			TraceID:       tx.TraceID,
			AccountID:     tx.AccountID,
			MerchantID:    tx.MerchantID,
			Amount:        tx.Amount,
			Status:        tx.Status,
			CreatedAt:     tx.CreatedAt,
			UpdatedAt:     tx.UpdatedAt,
		})
	}

	return items, nil
}
