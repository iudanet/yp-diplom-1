package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (r *postgresRepo) CreateOrder(ctx context.Context, userID int64, number string) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO orders (user_id, number) VALUES ($1, $2)",
		userID,
		number,
	)
	return err
}

func (r *postgresRepo) GetOrderByNumber(
	ctx context.Context,
	number string,
) (*models.OrderUser, error) {
	var order models.OrderUser
	var accrualCents sql.NullInt64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT number, status, accrual, uploaded_at
         FROM orders WHERE number = $1`,
		number,
	).Scan(&order.Number, &order.Status, &accrualCents, &order.UploadedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if accrualCents.Valid {
		order.Accrual = float64(accrualCents.Int64) / 100
		order.AccrualCents = accrualCents.Int64
	}

	return &order, nil
}
func (r *postgresRepo) GetOrdersByUserID(
	ctx context.Context,
	userID int64,
) ([]models.OrderUser, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT number, status, accrual, uploaded_at
         FROM orders WHERE user_id = $1
         ORDER BY uploaded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderUser
	for rows.Next() {
		var order models.OrderUser
		var accrualCents sql.NullInt64

		if err := rows.Scan(
			&order.Number,
			&order.Status,
			&accrualCents,
			&order.UploadedAt,
		); err != nil {
			return nil, err
		}

		if accrualCents.Valid {
			order.Accrual = float64(accrualCents.Int64) / 100
			order.AccrualCents = accrualCents.Int64
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *postgresRepo) GetOrderOwner(ctx context.Context, number string) (int64, error) {
	var userID int64
	err := r.db.QueryRowContext(
		ctx,
		"SELECT user_id FROM orders WHERE number = $1",
		number,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}
	return userID, nil
}
