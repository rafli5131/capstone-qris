package config

import (
	"time"

	"capstone-qris/internal/delivery/http"
	"capstone-qris/internal/delivery/http/middleware"
	"capstone-qris/internal/delivery/http/route"
	"capstone-qris/internal/repository"
	"capstone-qris/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	Redis     *redis.Client
	App       *fiber.App
	Log       *logrus.Logger
	Validator *validator.Validate
	Config    *Config
}

func Bootstrap(config *BootstrapConfig) {
	apiClientRepo := repository.NewApiClientRepository(config.DB, config.Log)
	merchantRepo := repository.NewMerchantRepository(config.DB, config.Log)
	accountRepo := repository.NewAccountRepository(config.DB, config.Log)
	transactionRepo := repository.NewTransactionRepository(config.DB, config.Log)

	qrisUseCase := usecase.NewQrisUseCase(
		config.DB, config.Redis, config.Log,
		time.Duration(config.Config.Redis.InquiryTTLSeconds)*time.Second,
		merchantRepo,
	)
	paymentUseCase := usecase.NewPaymentUseCase(
		config.DB, config.Redis, config.Log, config.Validator,
		accountRepo, transactionRepo, merchantRepo,
	)
	transactionUseCase := usecase.NewTransactionUseCase(
		config.DB, config.Log, transactionRepo, accountRepo,
	)

	sigMiddleware := middleware.NewSignatureMiddleware(
		config.DB, config.Log,
		apiClientRepo,
	)

	qrisController := http.NewQrisController(config.Log, qrisUseCase)
	paymentController := http.NewPaymentController(config.Log, paymentUseCase, transactionUseCase)

	routeConfig := &route.RouteConfig{
		App:                 config.App,
		QrisController:      qrisController,
		PaymentController:   paymentController,
		SignatureMiddleware: sigMiddleware,
	}
	routeConfig.Setup()
}
