package usecase

import (
	"context"
	"errors"
	"fmt"

	"capstone-qris/internal/entity"
	"capstone-qris/internal/model"
	"capstone-qris/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentUseCase struct {
	DB              *gorm.DB
	Redis           *redis.Client
	Log             *logrus.Logger
	Validator       *validator.Validate
	AccountRepo     *repository.AccountRepository
	TransactionRepo *repository.TransactionRepository
	MerchantRepo    *repository.MerchantRepository
}

func NewPaymentUseCase(
	db *gorm.DB,
	rdb *redis.Client,
	log *logrus.Logger,
	val *validator.Validate,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
	merchantRepo *repository.MerchantRepository,
) *PaymentUseCase {
	return &PaymentUseCase{
		DB:              db,
		Redis:           rdb,
		Log:             log,
		Validator:       val,
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
		MerchantRepo:    merchantRepo,
	}
}

func (uc *PaymentUseCase) Pay(ctx context.Context, accountID string, req *model.PaymentRequest) (*model.PaymentResponse, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	inquiryKey := fmt.Sprintf("qris:inquiry_id:%s", req.InquiryID)
	merchantID, err := uc.Redis.Get(ctx, inquiryKey).Result()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "inquiry_id is invalid or expired")
	}

	account, err := uc.AccountRepo.FindByAccountID(uc.DB, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "account not found")
		}
		uc.Log.WithError(err).Warn("failed to find account")
		return nil, fiber.ErrInternalServerError
	}

	if account.Balance < req.Amount {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "insufficient balance")
	}

	merchant, err := uc.MerchantRepo.FindByMerchantID(uc.DB, merchantID)
	if err != nil {
		uc.Log.WithError(err).Warnf("merchant not found: %s", merchantID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "merchant not found")
	}

	transactionID := uuid.New().String()
	traceID := "trace_" + uuid.New().String()[:8]

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	const maxRetries = 3
	var rowsAffected int64
	for i := 0; i < maxRetries; i++ {
		rowsAffected, err = uc.AccountRepo.UpdateWithOptimisticLock(tx, account.AccountID, req.Amount, account.Version)
		if err != nil {
			tx.Rollback()
			uc.Log.WithError(err).Warn("optimistic lock update failed")
			return nil, fiber.ErrInternalServerError
		}
		if rowsAffected > 0 {
			break
		}
		account, err = uc.AccountRepo.FindByAccountID(uc.DB, accountID)
		if err != nil || account.Balance < req.Amount {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "concurrent update conflict, please retry")
		}
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusConflict, "could not acquire lock after retries")
	}

	newTx := &entity.Transaction{
		TransactionID: transactionID,
		TraceID:       traceID,
		AccountID:     account.AccountID,
		MerchantID:    merchant.MerchantID,
		Amount:        req.Amount,
		Status:        "SUCCESS",
	}
	if err := uc.TransactionRepo.Create(tx, newTx); err != nil {
		tx.Rollback()
		uc.Log.WithError(err).Warn("failed to create transaction record")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Warn("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.PaymentResponse{
		Status:        "SUCCESS",
		TransactionID: transactionID,
		Message:       "Transaksi berhasil diproses",
	}, nil
}
