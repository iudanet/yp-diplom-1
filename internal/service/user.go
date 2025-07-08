package service

import "github.com/iudanet/yp-diplom-1/internal/repo"

func NewService(repo repo.Repositories) Service {
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

func (s *service) Login() error {
	return nil
}

func (s *service) Register() error {
	return nil
}
