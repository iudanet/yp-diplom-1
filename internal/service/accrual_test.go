package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccrualClient_GetOrderAccrual(t *testing.T) {
	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		orderNumber    string
		expectedResult *models.OrderAccrualResponse
		expectError    bool
	}{
		{
			name: "successful response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"order":"123","status":"PROCESSED","accrual":100.5}`))
			},
			orderNumber: "123",
			expectedResult: &models.OrderAccrualResponse{
				Order:   "123",
				Status:  "PROCESSED",
				Accrual: ptrFloat64(100.5),
			},
			expectError: false,
		},
		{
			name: "order not found",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			orderNumber:    "456",
			expectedResult: nil,
			expectError:    false,
		},
		{
			name: "too many requests",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Header().Set("Retry-After", "1")
			},
			orderNumber:    "789",
			expectedResult: nil,
			expectError:    true, // Ожидаем ошибку при 429
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewAccrualClient(server.URL)
			client.httpClient.Timeout = 100 * time.Millisecond

			result, err := client.GetOrderAccrual(context.Background(), tt.orderNumber)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func ptrFloat64(f float64) *float64 {
	return &f
}
