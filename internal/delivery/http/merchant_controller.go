package http

import (
	"capstone-qris/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MerchantController struct {
	Log             *logrus.Logger
	MerchantUseCase *usecase.MerchantUseCase
}

func NewMerchantController(log *logrus.Logger, merchantUseCase *usecase.MerchantUseCase) *MerchantController {
	return &MerchantController{Log: log, MerchantUseCase: merchantUseCase}
}

// GetMerchantIncome godoc
// @Summary      Merchant income overview
// @Description  Return total revenue and transactions for a merchant.
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        merchant_id  path      string  true  "Merchant ID"
// @Success      200  {object}  model.MerchantIncomeResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/merchant/{merchant_id}/income [get]
func (ctrl *MerchantController) GetMerchantIncome(c *fiber.Ctx) error {
	merchantID := c.Params("merchant_id")
	if merchantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "merchant_id is required")
	}

	resp, err := ctrl.MerchantUseCase.GetIncome(c.Context(), merchantID)
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to get merchant income")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// GetMerchantTransactions godoc
// @Summary      Merchant transactions
// @Description  Return all transactions for a merchant.
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        merchant_id  path      string  true  "Merchant ID"
// @Success      200  {object}  []model.TransactionListItem
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/merchant/{merchant_id}/transactions [get]
func (ctrl *MerchantController) GetMerchantTransactions(c *fiber.Ctx) error {
	merchantID := c.Params("merchant_id")
	if merchantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "merchant_id is required")
	}

	resp, err := ctrl.MerchantUseCase.GetTransactions(c.Context(), merchantID)
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to get merchant transactions")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}
