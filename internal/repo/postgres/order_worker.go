package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (r *postgresRepo) GetOrdersForProcessing(
	ctx context.Context,
	limit int,
) ([]models.OrderUser, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT number, status
		FROM orders
		WHERE status IN ($1, $2)
		ORDER BY created_at ASC
		LIMIT $3`,
		models.OrderUserStatusNew,
		models.OrderUserStatusProcessing,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderUser
	for rows.Next() {
		var order models.OrderUser
		if err := rows.Scan(&order.Number, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	// Проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *postgresRepo) UpdateOrderStatus(
	ctx context.Context,
	number string,
	status models.OrderUserStatus,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE orders
		SET status = $1, updated_at = $2
		WHERE number = $3`,
		status,
		time.Now(),
		number,
	)
	if err != nil {
		return err
	}

	// Проверяем, что строка была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when updating order status")
	}

	return nil
}

func (r *postgresRepo) UpdateOrderAccrual(
	ctx context.Context,
	number string,
	status models.OrderUserStatus,
	accrualCents int64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE orders
        SET status = $1, accrual = $2, updated_at = $3
        WHERE number = $4`,
		status,
		accrualCents,
		time.Now(),
		number,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when updating order accrual")
	}

	return nil
}
