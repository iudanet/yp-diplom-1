package models

import "time"

type OrderUserStatus = string

// - `NEW` — заказ загружен в систему, но не попал в обработку;
// - `PROCESSING` — вознаграждение за заказ рассчитывается;
// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.
const (
	OrderUserStatusNew        OrderUserStatus = "NEW"
	OrderUserStatusProcessing OrderUserStatus = "PROCESSING"
	OrderUserStatusInvalid    OrderUserStatus = "INVALID"
	OrderUserStatusProcessed  OrderUserStatus = "PROCESSED"
)

// GET /api/user/orders []OrderUser
type OrderUser struct {
	Number       string          `json:"number"`
	Status       OrderUserStatus `json:"status"`
	Accrual      float64         `json:"accrual,omitempty"`
	AccrualCents int64           `json:"-"` // копейки (только для внутреннего использования)
	UploadedAt   time.Time       `json:"uploaded_at"`
}
