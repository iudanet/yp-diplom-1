package service

import (
	"context"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

//
// * `POST /api/user/register` — регистрация пользователя;
// * `POST /api/user/login` — аутентификация пользователя;
// * `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
// * `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
// * `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
// * `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
// * `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем.

type Service interface {
	AuthService
	UserService
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (*models.UserAuth, error)
	Register(ctx context.Context, login, password string) error
}

type UserService interface {
	CreateOrders() error
	GetOrders() error
	GetBalance() error
	GetBalanceWithdrawals() error
	CreateWithdraw() error
}
