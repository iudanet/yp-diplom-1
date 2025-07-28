package repo

import (
	"context"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, login, passwordHash string) error
	GetUserByLogin(ctx context.Context, login string) (*models.UserAuth, error)
	GetUserByID(ctx context.Context, id int64) (*models.UserAuth, error)
}

// OrderRepository определяет методы для работы с заказами
type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int64, number string) error
	GetOrderByNumber(ctx context.Context, number string) (*models.OrderUser, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]models.OrderUser, error)
	GetOrderOwner(ctx context.Context, number string) (int64, error)
	GetOrdersForProcessing(ctx context.Context, limit int) ([]models.OrderUser, error)
	UpdateOrderStatus(ctx context.Context, number string, status models.OrderUserStatus) error
	UpdateOrderAccrual(
		ctx context.Context,
		number string,
		status models.OrderUserStatus,
		accrualCents int64,
	) error
}

// BalanceRepository определяет методы для работы с балансом и списаниями
type BalanceRepository interface {
	GetUserBalance(ctx context.Context, userID int64) (current, withdrawn int64, err error)
	CreateWithdrawal(ctx context.Context, userID int64, orderNumber string, sum int64) error
	GetWithdrawals(ctx context.Context, userID int64) ([]models.WithdrawalDB, error)
}

// Migrator определяет методы для миграций
type Migrator interface {
	Migrate(ctx context.Context) error
}

// Repositories объединяет все доменные репозитории
type Repositories interface {
	Migrator
	UserRepository
	OrderRepository
	BalanceRepository
}
