package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/repo/mock_repo"
	"github.com/iudanet/yp-diplom-1/internal/service/mock_service"
	"github.com/stretchr/testify/assert"
)

func TestWorkerProcessOrder(t *testing.T) {
	tests := []struct {
		name          string
		order         models.OrderUser
		accrualResult *models.OrderAccrualResponse
		accrualError  error
		mockSetup     func(*mock_repo.MockRepositories, *mock_service.MockAccrualClientInterface)
		wantErr       bool
	}{
		{
			name: "NewOrderUpdatesToProcessingAndThenInvalidWhenNotFound",
			order: models.OrderUser{
				Number: "123",
				Status: models.OrderUserStatusNew,
			},
			accrualResult: nil,
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				repoMock.EXPECT().
					UpdateOrderStatus(gomock.Any(), "123", models.OrderUserStatusProcessing).
					Return(nil)
				repoMock.EXPECT().
					UpdateOrderStatus(gomock.Any(), "123", models.OrderUserStatusInvalid).
					Return(nil)
				accrualMock.EXPECT().
					GetOrderAccrual(gomock.Any(), "123").
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "ProcessingOrderUpdatesToProcessedWithAccrual",
			order: models.OrderUser{
				Number: "456",
				Status: models.OrderUserStatusProcessing,
			},
			accrualResult: &models.OrderAccrualResponse{
				Order:   "456",
				Status:  models.OrderAccrualStatusProcessed,
				Accrual: 100.5,
			},
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				repoMock.EXPECT().UpdateOrderAccrual(
					gomock.Any(),
					"456",
					models.OrderUserStatusProcessed,
					int64(10050),
				).Return(nil)
				accrualMock.EXPECT().
					GetOrderAccrual(gomock.Any(), "456").
					Return(&models.OrderAccrualResponse{
						Order:   "456",
						Status:  models.OrderAccrualStatusProcessed,
						Accrual: 100.5,
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "ProcessingOrderUpdatesToInvalidWhenAccrualStatusInvalid",
			order: models.OrderUser{
				Number: "789",
				Status: models.OrderUserStatusProcessing,
			},
			accrualResult: &models.OrderAccrualResponse{
				Order:  "789",
				Status: models.OrderAccrualStatusInvalid,
			},
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				repoMock.EXPECT().UpdateOrderStatus(
					gomock.Any(),
					"789",
					models.OrderUserStatusInvalid,
				).Return(nil)
				accrualMock.EXPECT().
					GetOrderAccrual(gomock.Any(), "789").
					Return(&models.OrderAccrualResponse{
						Order:  "789",
						Status: models.OrderAccrualStatusInvalid,
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "ReturnsErrorWhenUpdateStatusFails",
			order: models.OrderUser{
				Number: "321",
				Status: models.OrderUserStatusNew,
			},
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				repoMock.EXPECT().UpdateOrderStatus(
					gomock.Any(),
					"321",
					models.OrderUserStatusProcessing,
				).Return(errors.New("update failed"))
				// Не ожидаем вызова GetOrderAccrual, так как предыдущая операция вернула ошибку
			},
			wantErr: true,
		},
		{
			name: "ProcessingOrderRemainsProcessingWhenAccrualStillProcessing",
			order: models.OrderUser{
				Number: "654",
				Status: models.OrderUserStatusProcessing,
			},
			accrualResult: &models.OrderAccrualResponse{
				Order:  "654",
				Status: models.OrderAccrualStatusProcessing,
			},
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				accrualMock.EXPECT().
					GetOrderAccrual(gomock.Any(), "654").
					Return(&models.OrderAccrualResponse{
						Order:  "654",
						Status: models.OrderAccrualStatusProcessing,
					}, nil)
				// No repo update expected as status remains the same
			},
			wantErr: false,
		},
		{
			name: "ReturnsErrorWhenAccrualClientFails",
			order: models.OrderUser{
				Number: "987",
				Status: models.OrderUserStatusProcessing,
			},
			accrualError: errors.New("accrual client error"),
			mockSetup: func(repoMock *mock_repo.MockRepositories, accrualMock *mock_service.MockAccrualClientInterface) {
				accrualMock.EXPECT().
					GetOrderAccrual(gomock.Any(), "987").
					Return(nil, errors.New("accrual client error"))
				// No repo calls expected as accrual client fails
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockRepositories(ctrl)
			mockAccrual := mock_service.NewMockAccrualClientInterface(ctrl)
			tt.mockSetup(mockRepo, mockAccrual)

			worker := NewWorker(mockRepo, mockAccrual)

			err := worker.processOrder(context.Background(), tt.order)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
