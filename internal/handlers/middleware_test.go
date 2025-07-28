package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/service/mock_service"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		mockSetup      func(*Server)
		expectedStatus int
		checkContext   bool
	}{
		{
			name: "ValidToken",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				token := generateTestToken(t, 123, "test_secret")
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			mockSetup: func(s *Server) {
				s.cfg.SecretKey = "test_secret"
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name: "MissingAuthorizationHeader",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			mockSetup:      func(s *Server) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "InvalidTokenFormat",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "InvalidTokenFormat")
				return req
			},
			mockSetup:      func(s *Server) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "InvalidToken",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				return req
			},
			mockSetup:      func(s *Server) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "ExpiredToken",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				// Генерация токена с истекшим сроком действия
				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": 123,
					"exp":     1, // В прошлом
				}).SignedString([]byte("test_secret"))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			mockSetup: func(s *Server) {
				s.cfg.SecretKey = "test_secret"
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_service.NewMockService(ctrl)
			cfg := config.New()
			srv := New(mockService, cfg)
			tt.mockSetup(srv)

			// Переменная для проверки вызова следующего обработчика
			var nextCalled bool
			var contextUserID interface{}

			// Тестовый обработчик
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				contextUserID = r.Context().Value(userIDKey)
				w.WriteHeader(http.StatusOK)
			})

			req := tt.setupRequest()
			w := httptest.NewRecorder()

			// Применяем middleware
			middleware := srv.authMiddleware(nextHandler)
			middleware.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.expectedStatus, resp.StatusCode, "неверный статус код")

			// Проверяем, был ли вызван следующий обработчик
			if tt.expectedStatus == http.StatusOK {
				require.True(t, nextCalled, "следующий обработчик не был вызван")
				if tt.checkContext {
					require.NotNil(t, contextUserID, "userID должен быть в контексте")
					require.Equal(
						t,
						int64(123),
						contextUserID.(int64),
						"неверный userID в контексте",
					)
				}
			} else {
				require.False(t, nextCalled, "следующий обработчик не должен быть вызван")
			}
		})
	}
}
