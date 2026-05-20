package middleware

import (
	"crypto/hmac"

	"capstone-qris/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SignatureMiddleware validates X-Client-Id and X-Client-Key on every request.
type SignatureMiddleware struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	ApiClientRepo *repository.ApiClientRepository
}

func NewSignatureMiddleware(
	db *gorm.DB,
	log *logrus.Logger,
	apiClientRepo *repository.ApiClientRepository,
) *SignatureMiddleware {
	return &SignatureMiddleware{
		DB:            db,
		Log:           log,
		ApiClientRepo: apiClientRepo,
	}
}

// Handle returns the Fiber handler function for client auth validation.
func (m *SignatureMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientID := c.Get("X-Client-Id")
		clientKey := c.Get("X-Client-Key")

		if clientID == "" || clientKey == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing required security headers")
		}

		client, err := m.ApiClientRepo.FindByClientID(m.DB, clientID)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unknown client id")
		}
		if client.Status != "ACTIVE" {
			return fiber.NewError(fiber.StatusUnauthorized, "client is not active")
		}

		if !hmac.Equal([]byte(clientKey), []byte(client.ClientSecret)) {
			m.Log.WithField("client_id", clientID).Warn("invalid client key")
			return fiber.NewError(fiber.StatusUnauthorized, "invalid client key")
		}

		c.Locals("client_id", clientID)
		return c.Next()
	}
}
