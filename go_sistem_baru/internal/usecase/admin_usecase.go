package usecase

import (
	"context"
	"errors"

	"capstone-qris/internal/entity"
	"capstone-qris/internal/model"
	"capstone-qris/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validator       *validator.Validate
	ApiClientRepo   *repository.ApiClientRepository
	TransactionRepo *repository.TransactionRepository
}

func NewAdminUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validator *validator.Validate,
	apiClientRepo *repository.ApiClientRepository,
	transactionRepo *repository.TransactionRepository,
) *AdminUseCase {
	return &AdminUseCase{
		DB:              db,
		Log:             log,
		Validator:       validator,
		ApiClientRepo:   apiClientRepo,
		TransactionRepo: transactionRepo,
	}
}

func (uc *AdminUseCase) ListTransactions(ctx context.Context) ([]model.TransactionListItem, error) {
	txs, err := uc.TransactionRepo.FindAll(uc.DB.WithContext(ctx))
	if err != nil {
		uc.Log.WithError(err).Warn("failed to list transactions")
		return nil, fiber.ErrInternalServerError
	}

	result := make([]model.TransactionListItem, 0, len(txs))
	for _, tx := range txs {
		result = append(result, model.TransactionListItem{
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

	return result, nil
}

func (uc *AdminUseCase) UpdateTransaction(ctx context.Context, transactionID string, req *model.TransactionUpdateRequest) (*entity.Transaction, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx, err := uc.TransactionRepo.FindByTransactionID(uc.DB.WithContext(ctx), transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "transaction not found")
		}
		uc.Log.WithError(err).Warn("failed to find transaction")
		return nil, fiber.ErrInternalServerError
	}

	if req.Amount > 0 {
		tx.Amount = req.Amount
	}
	if req.Status != "" {
		tx.Status = req.Status
	}

	if err := uc.TransactionRepo.Update(uc.DB.WithContext(ctx), tx); err != nil {
		uc.Log.WithError(err).Warn("failed to update transaction")
		return nil, fiber.ErrInternalServerError
	}

	return tx, nil
}

func (uc *AdminUseCase) ListApiClients(ctx context.Context) ([]model.ApiClientResponse, error) {
	clients, err := uc.ApiClientRepo.FindAll(uc.DB.WithContext(ctx))
	if err != nil {
		uc.Log.WithError(err).Warn("failed to list api clients")
		return nil, fiber.ErrInternalServerError
	}

	result := make([]model.ApiClientResponse, 0, len(clients))
	for _, client := range clients {
		result = append(result, model.ApiClientResponse{
			ClientID:  client.ClientID,
			Status:    client.Status,
			CreatedAt: client.CreatedAt,
		})
	}

	return result, nil
}

func (uc *AdminUseCase) CreateApiClient(ctx context.Context, req *model.ApiClientCreateRequest) (*model.ApiClientResponse, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := uc.ApiClientRepo.FindByClientID(uc.DB.WithContext(ctx), req.ClientID)
	if err == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "client id already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.Log.WithError(err).Warn("failed to check existing api client")
		return nil, fiber.ErrInternalServerError
	}

	client := &entity.ApiClient{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Status:       req.Status,
	}
	if err := uc.ApiClientRepo.Create(uc.DB.WithContext(ctx), client); err != nil {
		uc.Log.WithError(err).Warn("failed to create api client")
		return nil, fiber.ErrInternalServerError
	}

	return &model.ApiClientResponse{
		ClientID:  client.ClientID,
		Status:    client.Status,
		CreatedAt: client.CreatedAt,
	}, nil
}

func (uc *AdminUseCase) UpdateApiClient(ctx context.Context, clientID string, req *model.ApiClientUpdateRequest) (*model.ApiClientResponse, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	client, err := uc.ApiClientRepo.FindByClientID(uc.DB.WithContext(ctx), clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "client id not found")
		}
		uc.Log.WithError(err).Warn("failed to find api client")
		return nil, fiber.ErrInternalServerError
	}

	if req.ClientSecret != "" {
		client.ClientSecret = req.ClientSecret
	}
	if req.Status != "" {
		client.Status = req.Status
	}

	if err := uc.ApiClientRepo.Update(uc.DB.WithContext(ctx), client); err != nil {
		uc.Log.WithError(err).Warn("failed to update api client")
		return nil, fiber.ErrInternalServerError
	}

	return &model.ApiClientResponse{
		ClientID:  client.ClientID,
		Status:    client.Status,
		CreatedAt: client.CreatedAt,
	}, nil
}
