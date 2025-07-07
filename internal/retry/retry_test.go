package retry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type mockHTTPError struct {
	statusCode int
}

func (m *mockHTTPError) Error() string {
	return fmt.Sprintf("HTTP %d error", m.statusCode)
}

func (m *mockHTTPError) HTTPStatusCode() int {
	return m.statusCode
}

func TestIsRetriableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "HTTP500Error",
			err:      &mockHTTPError{statusCode: 500},
			expected: true,
		},
		{
			name:     "HTTP400Error",
			err:      &mockHTTPError{statusCode: 400},
			expected: false,
		},
		{
			name:     "networkError",
			err:      &net.DNSError{IsTimeout: true},
			expected: true,
		},
		{
			name:     "contextCanceled",
			err:      context.Canceled,
			expected: true,
		},
		{
			name:     "closedNetworkConnection",
			err:      errors.New("use of closed network connection"),
			expected: true,
		},
		{
			name:     "backendError",
			err:      errors.New("database is not available"),
			expected: false,
		},
		{
			name:     "non-retriableError",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "invalid_connection_error",
			err:      errors.New("invalid connection"),
			expected: false,
		},
		{
			name: "pqConnectionError",
			err: &pq.Error{
				Code: pgerrcode.ConnectionException,
			},
			expected: true,
		},
		{
			name: "pqTransactionResolutionUnknown",
			err: &pq.Error{
				Code: pgerrcode.TransactionResolutionUnknown,
			},
			expected: true,
		},
		{
			name:     "non-retriablePqError",
			err:      &pq.Error{Code: "23505"}, // unique_violation
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetriable(tt.err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
