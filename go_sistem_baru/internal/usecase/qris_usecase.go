package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"capstone-qris/internal/entity"
	"capstone-qris/internal/model"
	"capstone-qris/internal/model/converter"
	"capstone-qris/internal/repository"
	"capstone-qris/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// QrisUseCase handles QRIS inquiry business logic.
type QrisUseCase struct {
	DB           *gorm.DB
	Redis        *redis.Client
	Log          *logrus.Logger
	InquiryTTL   time.Duration
	MerchantRepo *repository.MerchantRepository
}

func NewQrisUseCase(
	db *gorm.DB,
	rdb *redis.Client,
	log *logrus.Logger,
	inquiryTTL time.Duration,
	merchantRepo *repository.MerchantRepository,
) *QrisUseCase {
	return &QrisUseCase{
		DB:           db,
		Redis:        rdb,
		Log:          log,
		InquiryTTL:   inquiryTTL,
		MerchantRepo: merchantRepo,
	}
}

// Inquiry returns merchant data for a given QRIS payload, using Redis as a cache.
func (uc *QrisUseCase) Inquiry(ctx context.Context, qrisPayload string) (*model.WebResponse[*model.InquiryResponse], error) {
	start := time.Now()

	// Build cache key from SHA-256 of the raw payload.
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(qrisPayload)))
	cacheKey := fmt.Sprintf("qris:inquiry:%s", hash)

	// --- Cache HIT ---
	cached, err := uc.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var inquiryResp model.InquiryResponse
		if jsonErr := json.Unmarshal([]byte(cached), &inquiryResp); jsonErr == nil {
			return &model.WebResponse[*model.InquiryResponse]{
				Status: "success",
				Data:   &inquiryResp,
				Metadata: &model.Metadata{
					LatencyMS: time.Since(start).Milliseconds(),
					Source:    "cache",
				},
			}, nil
		}
	}

	// --- Cache MISS ---
	parsed := util.ParseQRIS(qrisPayload)

	merchant, err := uc.MerchantRepo.FindByMerchantID(uc.DB, parsed.MerchantID)
	if err != nil {
		// Fallback: try to construct a merchant from parsed data if not found in DB.
		uc.Log.WithError(err).Warnf("merchant not found for id=%s, using parsed data", parsed.MerchantID)
		return nil, fmt.Errorf("merchant not found: %s", parsed.MerchantID)
	}

	// Generate inquiry_id and store inquiry data in Redis.
	inquiryID := "inq_" + uuid.New().String()[:8]
	ttl := uc.InquiryTTL

	respData := converter.ToInquiryResponse(merchant, inquiryID, parsed)

	// Cache the InquiryResponse so subsequent lookups with the same payload are fast.
	if jsonBytes, jsonErr := json.Marshal(respData); jsonErr == nil {
		uc.Redis.Set(ctx, cacheKey, string(jsonBytes), ttl)
	}

	// Also store inquiry_id → merchant_id mapping for the payment use case.
	inquiryKey := fmt.Sprintf("qris:inquiry_id:%s", inquiryID)
	uc.Redis.Set(ctx, inquiryKey, merchant.MerchantID, ttl)

	return &model.WebResponse[*model.InquiryResponse]{
		Status: "success",
		Data:   respData,
		Metadata: &model.Metadata{
			LatencyMS: time.Since(start).Milliseconds(),
			Source:    "database",
		},
	}, nil
}

// AddMerchantFromPayload creates or re-activates a merchant from a QRIS payload.
func (uc *QrisUseCase) AddMerchantFromPayload(ctx context.Context, qrisPayload string) (*model.WebResponse[*model.MerchantResponse], error) {
	parsed := util.ParseQRIS(qrisPayload)
	if parsed.MerchantID == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "merchant_id not found in qr payload")
	}
	if parsed.MerchantName == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "merchant_name not found in qr payload")
	}

	merchant, err := uc.MerchantRepo.FindByMerchantIDIncludingInactive(uc.DB, parsed.MerchantID)
	if err == nil {
		merchant.MerchantName = parsed.MerchantName
		if parsed.City != "" {
			merchant.City = parsed.City
		}
		if parsed.MCC != "" {
			merchant.MCC = parsed.MCC
		}
		if !merchant.IsActive {
			merchant.IsActive = true
		}
		if saveErr := uc.MerchantRepo.Save(uc.DB.WithContext(ctx), merchant); saveErr != nil {
			uc.Log.WithError(saveErr).Warn("failed to update merchant")
			return nil, fiber.ErrInternalServerError
		}
		return &model.WebResponse[*model.MerchantResponse]{
			Status: "success",
			Data: &model.MerchantResponse{
				MerchantID:   merchant.MerchantID,
				MerchantName: merchant.MerchantName,
				City:         merchant.City,
				MCC:          merchant.MCC,
				IsActive:     merchant.IsActive,
				IsNew:        false,
			},
		}, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.Log.WithError(err).Warn("failed to lookup merchant")
		return nil, fiber.ErrInternalServerError
	}

	newMerchant := &entity.Merchant{
		MerchantID:   parsed.MerchantID,
		MerchantName: parsed.MerchantName,
		City:         parsed.City,
		MCC:          parsed.MCC,
		IsActive:     true,
	}

	if createErr := uc.MerchantRepo.Create(uc.DB.WithContext(ctx), newMerchant); createErr != nil {
		uc.Log.WithError(createErr).Warn("failed to create merchant")
		return nil, fiber.ErrInternalServerError
	}

	return &model.WebResponse[*model.MerchantResponse]{
		Status: "success",
		Data: &model.MerchantResponse{
			MerchantID:   newMerchant.MerchantID,
			MerchantName: newMerchant.MerchantName,
			City:         newMerchant.City,
			MCC:          newMerchant.MCC,
			IsActive:     newMerchant.IsActive,
			IsNew:        true,
		},
	}, nil
}
