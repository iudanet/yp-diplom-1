package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (r *postgresRepo) CreateUser(ctx context.Context, login, passwordHash string) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users_auth (login, password_hash) VALUES ($1, $2)",
		login,
		passwordHash,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			return models.ErrUserAlreadyExists
		}
		return err
	}
	return nil
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
) (current, withdrawn int64, err error) {
	// Получаем сумму всех начислений в копейках (целое число)
	var accrued int64
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(accrual), 0)
         FROM orders
         WHERE user_id = $1 AND status = $2`,
		userID, models.OrderUserStatusProcessed,
	).Scan(&accrued)
	if err != nil {
		return 0, 0, err
	}

	// Получаем сумму всех списаний в копейках (целое число)
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(sum), 0)
         FROM withdrawals
         WHERE user_id = $1`,
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return 0, 0, err
	}

	// Текущий баланс = начисления - списания (все в копейках)
	current = accrued - withdrawn
	return current, withdrawn, nil
}

func (r *postgresRepo) CreateWithdrawal(
	ctx context.Context,
	userID int64,
	orderNumber string,
	sumCents int64,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	var accrued int64
	err = tx.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(accrual), 0)
         FROM orders
         WHERE user_id = $1 AND status = $2`,
		userID, models.OrderUserStatusProcessed,
	).Scan(&accrued)
	if err != nil {
		return fmt.Errorf("failed to get accrued sum: %w", err)
	}

	var withdrawn int64
	err = tx.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(sum), 0)
         FROM withdrawals
         WHERE user_id = $1`,
		userID,
	).Scan(&withdrawn)
	if err != nil {
		return fmt.Errorf("failed to get withdrawn sum: %w", err)
	}

	currentBalance := accrued - withdrawn
	if currentBalance < sumCents {
		return models.ErrInsufficientFunds
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, sum)
         VALUES ($1, $2, $3)`,
		userID, orderNumber, sumCents,
	)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepo) GetWithdrawals(
	ctx context.Context,
	userID int64,
) ([]models.WithdrawalDB, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT order_number, sum, processed_at
         FROM withdrawals
         WHERE user_id = $1
         ORDER BY processed_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.WithdrawalDB
	for rows.Next() {
		var w models.WithdrawalDB
		if err := rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
