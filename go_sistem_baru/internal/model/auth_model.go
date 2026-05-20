package model

type RegisterRequest struct {
	Username       string  `json:"username" validate:"required,alphanum"`
	Password       string  `json:"password" validate:"required,min=8"`
	InitialBalance float64 `json:"initial_balance" validate:"required,gte=0"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token     string  `json:"token"`
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type AccountProfileResponse struct {
	AccountID string  `json:"account_id"`
	Username  string  `json:"username"`
	Balance   float64 `json:"balance"`
}
