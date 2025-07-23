package models

// `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
// `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
// `PROCESSING` — расчёт начисления в процессе;
// `PROCESSED` — расчёт начисления окончен;
type OrderAccrualStatus = string

const (
	OrderAccrualStatusProcessing OrderAccrualStatus = "PROCESSING"
	OrderAccrualStatusInvalid    OrderAccrualStatus = "INVALID"
	OrderAccrualStatusProcessed  OrderAccrualStatus = "PROCESSED"
	OrderAccrualStatusRegistered OrderAccrualStatus = "REGISTERED"
)

// GET /api/orders/{number} HTTP/1.1
type OrderAccrualResponse struct {
	Order   string             `json:"order"`
	Status  OrderAccrualStatus `json:"status"`
	Accrual float64            `json:"accrual,omitempty"` // Изменено с int на float64
}
