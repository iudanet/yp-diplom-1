package models

import "time"

// снятие денег со счета
// GET /api/user/withdrawals HTTP/1.1 []Withdrawal
type Withdrawal struct {
	Order       string    `json:"-"`
	Sum         int       `json:"-"` // копейки
	ProcessedAt time.Time `json:"-"`
}

// Для хранения в БД
type WithdrawalDB struct {
	Order       string
	Sum         int // копейки
	ProcessedAt time.Time
}
