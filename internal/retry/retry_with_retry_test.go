package retry

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestWithRetry(t *testing.T) {
	tests := []struct {
		name      string
		op        func() error
		wantError bool
	}{
		{
			name: "successOnFirstAttempt",
			op: func() error {
				return nil
			},
			wantError: false,
		},
		{
			name: "retriableErrorThenSuccess",
			op: func() func() error {
				attempt := 0
				return func() error {
					if attempt < 2 {
						attempt++
						return &pq.Error{Code: pgerrcode.ConnectionException}
					}
					return nil
				}
			}(),
			wantError: false,
		},
		{
			name: "nonRetriableError",
			op: func() error {
				return errors.New("non-retriable error")
			},
			wantError: true,
		},
		{
			name: "allAttemptsFail",
			op: func() error {
				return &pq.Error{Code: pgerrcode.ConnectionException}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WithRetry(tt.op)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
