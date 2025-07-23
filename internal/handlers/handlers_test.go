package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/service/mock_service"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestServer_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        models.RegisterRequest
		mockSetup      func(*mock_service.MockService)
		expectedStatus int
	}{
		{
			name: "successful registration",
			request: models.RegisterRequest{
				Login:    "testuser",
				Password: "testpass",
			},
			mockSetup: func(mock *mock_service.MockService) {
				mock.EXPECT().Register(gomock.Any(), "testuser", "testpass").Return(nil)
				mock.EXPECT().
					Login(gomock.Any(), "testuser", "testpass").
					Return(&models.UserAuth{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty login or password",
			request: models.RegisterRequest{
				Login:    "",
				Password: "",
			},
			mockSetup:      func(mock *mock_service.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "user_already_exists",
			request: models.RegisterRequest{
				Login:    "existing",
				Password: "testpass",
			},
			mockSetup: func(mock *mock_service.MockService) {
				mock.EXPECT().Register(gomock.Any(), "existing", "testpass").
					Return(&pq.Error{Code: pgerrcode.UniqueViolation})
				// Имитируем ошибку PostgreSQL
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockService(ctrl)
			tt.mockSetup(mockService)

			cfg := config.New()
			srv := New(mockService, cfg)

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
			w := httptest.NewRecorder()

			srv.Register(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestServer_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        models.LoginRequest
		mockSetup      func(*mock_service.MockService)
		expectedStatus int
	}{
		{
			name: "successful login",
			request: models.LoginRequest{
				Login:    "testuser",
				Password: "testpass",
			},
			mockSetup: func(mock *mock_service.MockService) {
				mock.EXPECT().
					Login(gomock.Any(), "testuser", "testpass").
					Return(&models.UserAuth{ID: 1}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			request: models.LoginRequest{
				Login:    "testuser",
				Password: "wrongpass",
			},
			mockSetup: func(mock *mock_service.MockService) {
				mock.EXPECT().
					Login(gomock.Any(), "testuser", "wrongpass").
					Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockService(ctrl)
			tt.mockSetup(mockService)

			cfg := config.New()
			srv := New(mockService, cfg)

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			srv.Login(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestServer_PostOrders(t *testing.T) {
	tests := []struct {
		name           string
		orderNumber    string
		mockSetup      func(*mock_service.MockService)
		expectedStatus int
	}{
		{
			name:        "successful order upload",
			orderNumber: "4561261212345467", // Valid Luhn number
			mockSetup: func(mock *mock_service.MockService) {
				mock.EXPECT().CreateOrder(gomock.Any(), int64(1), "4561261212345467").Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "invalid order number",
			orderNumber:    "123",
			mockSetup:      func(mock *mock_service.MockService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockService(ctrl)
			tt.mockSetup(mockService)

			cfg := config.New()
			srv := New(mockService, cfg)

			req := httptest.NewRequest(
				"POST",
				"/api/user/orders",
				bytes.NewReader([]byte(tt.orderNumber)),
			)
			req.Header.Set("Content-Type", "text/plain")
			req.Header.Set("Authorization", "Bearer "+generateTestToken(t, 1, cfg.SecretKey))

			w := httptest.NewRecorder()

			srv.PostOrders(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func generateTestToken(t *testing.T, userID int64, secretKey string) string {
	token, err := generateToken(userID, secretKey)
	require.NoError(t, err)
	return token
}
