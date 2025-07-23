package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (s *service) CreateOrder(ctx context.Context, userID int64, number string) error {
	if !isValidLuhn(number) {
		return models.ErrInvalidOrderNumber
	}

	existingOrder, err := s.repo.GetOrderByNumber(ctx, number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check order: %w", err)
	}

	if existingOrder != nil {
		ownerID, err := s.repo.GetOrderOwner(ctx, number)
		if err != nil {
			return fmt.Errorf("failed to get order owner: %w", err)
		}

		if ownerID == userID {
			return models.ErrOrderAlreadyUploaded
		}
		return models.ErrOrderAlreadyUploadedByAnotherUser
	}

	return s.repo.CreateOrder(ctx, userID, number)
}

func (s *service) GetOrders(ctx context.Context, userID int64) ([]models.OrderUser, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

// isValidLuhn проверяет номер заказа по алгоритму Луна
func isValidLuhn(number string) bool {
	sum := 0
	nDigits := len(number)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit := int(number[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
