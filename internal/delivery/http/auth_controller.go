package http

import (
	"capstone-qris/internal/model"
	"capstone-qris/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(log *logrus.Logger, authUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{Log: log, AuthUseCase: authUseCase}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new account with initial balance and return JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.RegisterRequest  true  "Register request"
// @Success      201  {object}  model.AuthResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /api/auth/register [post]
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	resp, err := ctrl.AuthUseCase.Register(c.Context(), &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("register failed")
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": resp})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate with username and password and return JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.LoginRequest  true  "Login request"
// @Success      200  {object}  model.AuthResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /api/auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	resp, err := ctrl.AuthUseCase.Login(c.Context(), &req)
	if err != nil {
		ctrl.Log.WithError(err).Warn("login failed")
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": resp})
}
