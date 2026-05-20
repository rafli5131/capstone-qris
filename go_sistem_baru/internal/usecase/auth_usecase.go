package usecase

import (
	"context"
	"errors"
	"time"

	"capstone-qris/internal/entity"
	"capstone-qris/internal/model"
	"capstone-qris/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Validator     *validator.Validate
	AccountRepo   *repository.AccountRepository
	JwtSecret     string
	TokenDuration time.Duration
}

func NewAuthUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validator *validator.Validate,
	accountRepo *repository.AccountRepository,
	jwtSecret string,
	tokenDuration time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		DB:            db,
		Log:           log,
		Validator:     validator,
		AccountRepo:   accountRepo,
		JwtSecret:     jwtSecret,
		TokenDuration: tokenDuration,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := uc.AccountRepo.FindByUsername(uc.DB.WithContext(ctx), req.Username)
	if err == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "username already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.Log.WithError(err).Warn("failed to check username")
		return nil, fiber.ErrInternalServerError
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	account := &entity.Account{
		AccountID:    "acct_" + uuid.New().String()[:8],
		Username:     req.Username,
		PasswordHash: string(hash),
		Balance:      req.InitialBalance,
	}

	if err := uc.AccountRepo.Create(uc.DB.WithContext(ctx), account); err != nil {
		uc.Log.WithError(err).Warn("failed to create account")
		return nil, fiber.ErrInternalServerError
	}

	token, err := uc.generateToken(account.AccountID)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to generate token")
		return nil, fiber.ErrInternalServerError
	}

	return &model.AuthResponse{Token: token, AccountID: account.AccountID, Balance: account.Balance}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	if err := uc.Validator.StructCtx(ctx, req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	account, err := uc.AccountRepo.FindByUsername(uc.DB.WithContext(ctx), req.Username)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := uc.generateToken(account.AccountID)
	if err != nil {
		uc.Log.WithError(err).Warn("failed to generate token")
		return nil, fiber.ErrInternalServerError
	}

	return &model.AuthResponse{Token: token, AccountID: account.AccountID, Balance: account.Balance}, nil
}

func (uc *AuthUseCase) generateToken(accountID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   accountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(uc.TokenDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(uc.JwtSecret))
}
