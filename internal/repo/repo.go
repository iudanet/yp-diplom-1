package repo

import (
	"context"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

type Repositories interface {
	Migrate(ctx context.Context) error
	CreateUser(ctx context.Context, login, passwordHash string) error
	GetUserByLogin(ctx context.Context, login string) (*models.UserAuth, error)
}
