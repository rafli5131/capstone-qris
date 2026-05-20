package route

import (
	deliveryHttp "capstone-qris/internal/delivery/http"
	"capstone-qris/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	_ "capstone-qris/docs"
)

type RouteConfig struct {
	App                *fiber.App
	QrisController     *deliveryHttp.QrisController
	PaymentController  *deliveryHttp.PaymentController
	MerchantController *deliveryHttp.MerchantController
	AdminController    *deliveryHttp.AdminController
	AuthController     *deliveryHttp.AuthController
	JWTMiddleware      *middleware.JWTMiddleware
}

func (r *RouteConfig) Setup() {
	r.App.Get("/swagger/*", fiberSwagger.HandlerDefault)
	r.App.Get("/swagger", fiberSwagger.HandlerDefault)

	r.App.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := r.App.Group("/api/auth")
	auth.Post("/register", r.AuthController.Register)
	auth.Post("/login", r.AuthController.Login)

	api := r.App.Group("/api", r.JWTMiddleware.Handle())

	qris := api.Group("/qris")
	qris.Get("/inquiry/:qris_payload", r.QrisController.Inquiry)
	qris.Post("/inquiry/image", r.QrisController.InquiryFromImage)
	qris.Post("/merchant/image", r.QrisController.AddMerchantFromImage)
	qris.Post("/payment", r.PaymentController.Pay)
	qris.Get("/status/:transaction_id", r.PaymentController.GetTransactionStatus)

	admin := api.Group("/admin")
	admin.Get("/transactions", r.AdminController.ListTransactions)
	admin.Put("/transactions/:transaction_id", r.AdminController.UpdateTransaction)
	admin.Get("/api-clients", r.AdminController.ListApiClients)
	admin.Post("/api-clients", r.AdminController.CreateApiClient)
	admin.Put("/api-clients/:client_id", r.AdminController.UpdateApiClient)

	merchant := api.Group("/merchant")
	merchant.Get(":merchant_id/income", r.MerchantController.GetMerchantIncome)
	merchant.Get(":merchant_id/transactions", r.MerchantController.GetMerchantTransactions)
}
