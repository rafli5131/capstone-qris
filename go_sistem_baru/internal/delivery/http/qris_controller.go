package http

import (
	"io"
	"net/url"

	"capstone-qris/internal/usecase"
	"capstone-qris/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type QrisController struct {
	Log         *logrus.Logger
	QrisUseCase *usecase.QrisUseCase
}

func NewQrisController(log *logrus.Logger, qrisUseCase *usecase.QrisUseCase) *QrisController {
	return &QrisController{Log: log, QrisUseCase: qrisUseCase}
}

// Inquiry godoc
// @Summary      QRIS inquiry
// @Description  Get merchant data by QRIS payload (decoded if URL-encoded).
// @Tags         qris
// @Accept       json
// @Produce      json
// @Param        qris_payload  path      string  true  "QRIS payload (raw or URL-encoded)"
// @Success      200  {object}  model.InquiryResponseEnvelope
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/qris/inquiry/{qris_payload} [get]
func (ctrl *QrisController) Inquiry(c *fiber.Ctx) error {
	rawPayload := c.Params("qris_payload")

	decoded, err := url.QueryUnescape(rawPayload)
	if err != nil {
		decoded = rawPayload
	}

	resp, err := ctrl.QrisUseCase.Inquiry(c.Context(), decoded)
	if err != nil {
		ctrl.Log.WithError(err).Warn("inquiry failed")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// InquiryFromImage godoc
// @Summary      QRIS inquiry from image
// @Description  Get merchant data by QRIS image (QR code).
// @Tags         qris
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "QR code image"
// @Success      200  {object}  model.InquiryResponseEnvelope
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/qris/inquiry/image [post]
func (ctrl *QrisController) InquiryFromImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "image file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to open image file")
	}
	defer file.Close()

	payload, err := util.DecodeQRPayloadFromImage(io.LimitReader(file, 10<<20))
	if err != nil {
		ctrl.Log.WithError(err).Warn("decode qr image failed")
		return fiber.NewError(fiber.StatusBadRequest, "invalid qr image")
	}

	resp, err := ctrl.QrisUseCase.Inquiry(c.Context(), payload)
	if err != nil {
		ctrl.Log.WithError(err).Warn("inquiry from image failed")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// AddMerchantFromImage godoc
// @Summary      Add merchant from image
// @Description  Create or re-activate a merchant using a QRIS image (QR code).
// @Tags         merchant
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "QR code image"
// @Success      200  {object}  model.MerchantResponseEnvelope
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /api/qris/merchant/image [post]
func (ctrl *QrisController) AddMerchantFromImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "image file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to open image file")
	}
	defer file.Close()

	payload, err := util.DecodeQRPayloadFromImage(io.LimitReader(file, 10<<20))
	if err != nil {
		ctrl.Log.WithError(err).Warn("decode qr image failed")
		return fiber.NewError(fiber.StatusBadRequest, "invalid qr image")
	}

	resp, err := ctrl.QrisUseCase.AddMerchantFromPayload(c.Context(), payload)
	if err != nil {
		ctrl.Log.WithError(err).Warn("add merchant from image failed")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
