package models

import "time"

// снятие денег со счета
// GET /api/user/withdrawals HTTP/1.1 []Withdrawal
type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
