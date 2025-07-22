package repo

import (
	"context"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

type Repositories interface {
	Migrate(ctx context.Context) error
	CreateUser(ctx context.Context, login, passwordHash string) error
	GetUserByLogin(ctx context.Context, login string) (*models.UserAuth, error)
	GetUserByID(ctx context.Context, id int64) (*models.UserAuth, error)
	GetUserBalance(ctx context.Context, userID int64) (current, withdrawn float64, err error)
	// Методы для работы с заказами
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
		accrual float64,
	) error
}
