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

func (s *service) CreateOrders() error {
	return nil
}

func (s *service) GetOrders() error {
	return nil
}

func (s *service) GetBalance() error {
	return nil
}

func (s *service) GetBalanceWithdrawals() error {
	return nil
}

func (s *service) CreateWithdraw() error {
	return nil
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
