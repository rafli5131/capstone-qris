package http

import (
	"capstone-qris/internal/model"
	"capstone-qris/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AdminController struct {
	Log          *logrus.Logger
	AdminUseCase *usecase.AdminUseCase
}

func NewAdminController(log *logrus.Logger, adminUseCase *usecase.AdminUseCase) *AdminController {
	return &AdminController{Log: log, AdminUseCase: adminUseCase}
}

// ListTransactions godoc
// @Summary      List all transactions
// @Description  Return all transaction records for admin review.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.TransactionListItem
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/admin/transactions [get]
func (ctrl *AdminController) ListTransactions(c *fiber.Ctx) error {
	list, err := ctrl.AdminUseCase.ListTransactions(c.Context())
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to list transactions")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}

// UpdateTransaction godoc
// @Summary      Update transaction
// @Description  Modify transaction status or amount.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        transaction_id  path      string                        true  "Transaction ID"
// @Param        request         body      model.TransactionUpdateRequest true  "Transaction update request"
// @Success      200  {object}  model.TransactionListItem
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/admin/transactions/{transaction_id} [put]
func (ctrl *AdminController) UpdateTransaction(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")
	if transactionID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "transaction_id is required")
	}

	var req model.TransactionUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := ctrl.AdminUseCase.UpdateTransaction(c.Context(), transactionID, &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to update transaction")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": model.TransactionListItem{
			TransactionID: updated.TransactionID,
			TraceID:       updated.TraceID,
			AccountID:     updated.AccountID,
			MerchantID:    updated.MerchantID,
			Amount:        updated.Amount,
			Status:        updated.Status,
			CreatedAt:     updated.CreatedAt,
			UpdatedAt:     updated.UpdatedAt,
		},
	})
}

// ListApiClients godoc
// @Summary      List API clients
// @Description  Return all configured API clients.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.ApiClientResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/admin/api-clients [get]
func (ctrl *AdminController) ListApiClients(c *fiber.Ctx) error {
	clients, err := ctrl.AdminUseCase.ListApiClients(c.Context())
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to list api clients")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   clients,
	})
}

// CreateApiClient godoc
// @Summary      Create API client
// @Description  Add a new API client with client id and client key.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body      model.ApiClientCreateRequest  true  "API client create request"
// @Success      201  {object}  model.ApiClientResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/admin/api-clients [post]
func (ctrl *AdminController) CreateApiClient(c *fiber.Ctx) error {
	var req model.ApiClientCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	client, err := ctrl.AdminUseCase.CreateApiClient(c.Context(), &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to create api client")
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   client,
	})
}

// UpdateApiClient godoc
// @Summary      Update API client
// @Description  Edit an existing client secret or status.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        client_id  path      string                      true  "Client ID"
// @Param        request    body      model.ApiClientUpdateRequest true  "API client update request"
// @Success      200  {object}  model.ApiClientResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/admin/api-clients/{client_id} [put]
func (ctrl *AdminController) UpdateApiClient(c *fiber.Ctx) error {
	clientID := c.Params("client_id")
	if clientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "client_id is required")
	}

	var req model.ApiClientUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	client, err := ctrl.AdminUseCase.UpdateApiClient(c.Context(), clientID, &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("failed to update api client")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   client,
	})
}
