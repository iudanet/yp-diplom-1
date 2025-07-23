package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *postgresRepo) CreateWithdrawal(ctx context.Context, userID int64, orderNumber string, sum float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Получаем сумму начислений
	var accrued float64
	err = tx.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(accrual), 0)
         FROM orders
         WHERE user_id = $1 AND status = $2`,
		userID, models.OrderUserStatusProcessed,
	).Scan(&accrued)
	if err != nil {
		return fmt.Errorf("failed to get accrued sum: %w", err)
	}

	// Получаем сумму списаний
	var withdrawn float64
	err = tx.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(sum), 0)
         FROM withdrawals
         WHERE user_id = $1`,
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return fmt.Errorf("failed to get withdrawn sum: %w", err)
	}

	// Проверяем баланс
	currentBalance := accrued - withdrawn
	if currentBalance < sum {
		return models.ErrInsufficientFunds
	}

	// Создаем списание
	_, err = tx.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, sum)
         VALUES ($1, $2, $3)`,
		userID, orderNumber, sum,
	)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal: %w", err)
	}

	return tx.Commit()
}
