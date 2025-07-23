package models

import "time"

// POST /api/user/login HTTP/1.1
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// POST /api/user/register HTTP/1.1
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// POST /api/user/balance/withdraw HTTP/1.1
type BalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// GET /api/user/balance HTTP/1.1
type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"` // рубли
	ProcessedAt time.Time `json:"processed_at"`
}
