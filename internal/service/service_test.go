package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/repo/mock_repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		password  string
		mockSetup func(*mock_repo.MockRepositories)
		wantErr   bool
	}{
		{
			name:     "successful registration",
			login:    "newuser",
			password: "password",
			mockSetup: func(mock *mock_repo.MockRepositories) {
				mock.EXPECT().
					GetUserByLogin(gomock.Any(), "newuser").
					Return(nil, errors.New("not found"))
				mock.EXPECT().CreateUser(gomock.Any(), "newuser", gomock.Any()).DoAndReturn(
					func(ctx context.Context, login, hash string) error {
						err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("password"))
						assert.NoError(t, err)
						return nil
					},
				)
			},
			wantErr: false,
		},
		{
			name:     "user already exists",
			login:    "existing",
			password: "password",
			mockSetup: func(mock *mock_repo.MockRepositories) {
				mock.EXPECT().
					GetUserByLogin(gomock.Any(), "existing").
					Return(&models.UserAuth{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockRepositories(ctrl)
			tt.mockSetup(mockRepo)

			s := New(mockRepo)

			err := s.Register(context.Background(), tt.login, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	tests := []struct {
		name      string
		login     string
		password  string
		mockSetup func(*mock_repo.MockRepositories)
		wantErr   bool
	}{
		{
			name:     "successful login",
			login:    "user",
			password: "correct",
			mockSetup: func(mock *mock_repo.MockRepositories) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
				mock.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(&models.UserAuth{
					PasswordHash: string(hash),
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "wrong password",
			login:    "user",
			password: "wrong",
			mockSetup: func(mock *mock_repo.MockRepositories) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
				mock.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(&models.UserAuth{
					PasswordHash: string(hash),
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockRepositories(ctrl)
			tt.mockSetup(mockRepo)

			s := New(mockRepo)

			_, err := s.Login(context.Background(), tt.login, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_CreateOrder(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		number    string
		mockSetup func(*mock_repo.MockRepositories)
		wantErr   error
	}{
		{
			name:   "successful order creation",
			userID: 1,
			number: "4561261212345467", // Valid Luhn
			mockSetup: func(mock *mock_repo.MockRepositories) {
				mock.EXPECT().GetOrderByNumber(gomock.Any(), "4561261212345467").Return(nil, nil)
				mock.EXPECT().CreateOrder(gomock.Any(), int64(1), "4561261212345467").Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "invalid order number",
			userID: 1,
			number: "123",
			mockSetup: func(mock *mock_repo.MockRepositories) {
				// No expectations as validation should fail before repo calls
			},
			wantErr: models.ErrInvalidOrderNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockRepositories(ctrl)
			tt.mockSetup(mockRepo)

			s := New(mockRepo)

			err := s.CreateOrder(context.Background(), tt.userID, tt.number)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
