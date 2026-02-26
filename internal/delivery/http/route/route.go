package route

import (
	deliveryHttp "capstone-qris/internal/delivery/http"
	"capstone-qris/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	_ "capstone-qris/docs"
)

type RouteConfig struct {
	App                 *fiber.App
	QrisController      *deliveryHttp.QrisController
	PaymentController   *deliveryHttp.PaymentController
	SignatureMiddleware *middleware.SignatureMiddleware
}

func (r *RouteConfig) Setup() {
	r.App.Get("/swagger/*", fiberSwagger.HandlerDefault)

	r.App.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := r.App.Group("/api", r.SignatureMiddleware.Handle())

	qris := api.Group("/qris")
	qris.Get("/inquiry/:qris_payload", r.QrisController.Inquiry)
	qris.Post("/inquiry/image", r.QrisController.InquiryFromImage)
	qris.Post("/merchant/image", r.QrisController.AddMerchantFromImage)
	qris.Post("/payment", r.PaymentController.Pay)
	qris.Get("/status/:transaction_id", r.PaymentController.GetTransactionStatus)
}
