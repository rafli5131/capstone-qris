package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewFiber(cfg *Config) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      cfg.Web.Prefork,
		ErrorHandler: NewErrorHandler(),
		AppName:      "QRIS Payment API",
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(healthcheck.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		return ctx.Status(code).JSON(fiber.Map{
			"status": "error",
			"errors": err.Error(),
		})
	}
}
