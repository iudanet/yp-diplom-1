package service

import (
	"context"
	"fmt"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

func New(repo repo.Repositories) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo repo.Repositories
}

func (s *service) Login(ctx context.Context, login, password string) (*models.UserAuth, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Register(ctx context.Context, login, password string) error {
	// Check if user already exists
	_, err := s.repo.GetUserByLogin(ctx, login)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.CreateUser(ctx, login, string(passwordHash))
}

func (s *service) GetUserBalance(
	ctx context.Context,
	userID int64,
) (current, withdrawn float64, err error) {
	currentCents, withdrawnCents, err := s.repo.GetUserBalance(ctx, userID)

	// Конвертируем копейки в рубли с точностью до 2 знаков
	currentRub := roundToTwoDecimals(float64(currentCents) / 100)
	withdrawnRub := roundToTwoDecimals(float64(withdrawnCents) / 100)
	return currentRub, withdrawnRub, err
}

func (s *service) CreateWithdrawal(
	ctx context.Context,
	userID int64,
	orderNumber string,
	sumCents int64,
) error {
	return s.repo.CreateWithdrawal(ctx, userID, orderNumber, sumCents)
}

func (s *service) GetWithdrawals(ctx context.Context, userID int64) ([]models.WithdrawalDB, error) {
	return s.repo.GetWithdrawals(ctx, userID)
}
