package http

import (
	"capstone-qris/internal/model"
	"capstone-qris/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PaymentController struct {
	Log                *logrus.Logger
	PaymentUseCase     *usecase.PaymentUseCase
	TransactionUseCase *usecase.TransactionUseCase
}

func NewPaymentController(
	log *logrus.Logger,
	paymentUseCase *usecase.PaymentUseCase,
	transactionUseCase *usecase.TransactionUseCase,
) *PaymentController {
	return &PaymentController{
		Log:                log,
		PaymentUseCase:     paymentUseCase,
		TransactionUseCase: transactionUseCase,
	}
}

// Pay godoc
// @Summary      Create payment
// @Description  Submit payment for a QRIS inquiry.
// @Tags         qris
// @Accept       json
// @Produce      json
// @Param        request  body      model.PaymentRequest  true  "Payment request"
// @Success      202  {object}  model.PaymentResponseEnvelope
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      422  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     X-Client-Id
// @Security     X-Client-Key
// @Router       /api/qris/payment [post]
func (ctrl *PaymentController) Pay(c *fiber.Ctx) error {
	var req model.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	resp, err := ctrl.PaymentUseCase.Pay(c.Context(), &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("payment failed")
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// GetTransactionStatus godoc
// @Summary      Get transaction status
// @Description  Fetch current status and final balance for a transaction.
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        transaction_id  path      string  true  "Transaction ID"
// @Success      200  {object}  model.TransactionStatusResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     X-Client-Id
// @Security     X-Client-Key
// @Router       /api/qris/status/{transaction_id} [get]
func (ctrl *PaymentController) GetTransactionStatus(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")
	if transactionID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "transaction_id is required")
	}

	resp, err := ctrl.TransactionUseCase.GetStatus(c.Context(), transactionID)
	if err != nil {
		ctrl.Log.WithError(err).Warn("get transaction status failed")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
