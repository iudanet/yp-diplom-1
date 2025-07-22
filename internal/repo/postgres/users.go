package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (r *postgresRepo) CreateUser(ctx context.Context, login, passwordHash string) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users_auth (login, password_hash) VALUES ($1, $2)",
		login,
		passwordHash,
	)
	return err
}

func (r *postgresRepo) GetUserByLogin(ctx context.Context, login string) (*models.UserAuth, error) {
	var user models.UserAuth
	err := r.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users_auth WHERE login = $1", login).
		Scan(&user.ID, &user.Login, &user.PasswordHash)
	return &user, err
}

func (r *postgresRepo) GetUserByID(ctx context.Context, id int64) (*models.UserAuth, error) {
	var user models.UserAuth
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, login, password_hash FROM users_auth WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) GetUserBalance(
	ctx context.Context,
	userID int64,
) (current, withdrawn float64, err error) {
	// Получаем текущий баланс (сумма всех начислений)
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(accrual), 0)
         FROM orders
         WHERE user_id = $1 AND status = $2`,
		userID, models.OrderUserStatusProcessed,
	).Scan(&current)
	if err != nil {
		return 0, 0, err
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(sum), 0)
	     FROM withdrawals
	     WHERE user_id = $1`,
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return 0, 0, err
	}

	return current, withdrawn, nil
}
