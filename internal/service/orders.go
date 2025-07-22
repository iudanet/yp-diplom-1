package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (s *service) CreateOrder(ctx context.Context, userID int64, number string) error {
	// Проверяем валидность номера заказа
	if !isValidLuhn(number) {
		return models.ErrInvalidOrderNumber
	}

	// Проверяем существование пользователя
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrUserNotFound
		}
		return err
	}

	// Проверяем, существует ли уже такой заказ
	existingOrder, err := s.repo.GetOrderByNumber(ctx, number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingOrder != nil {
		// Если заказ уже существует, проверяем принадлежность пользователю
		orderOwner, err := s.repo.GetOrderOwner(ctx, number)
		if err != nil {
			return err
		}

		if orderOwner == userID {
			return models.ErrOrderAlreadyUploaded
		}
		return models.ErrOrderAlreadyUploadedByAnotherUser
	}

	// Создаем новый заказ
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
