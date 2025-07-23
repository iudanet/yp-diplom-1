package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundToTwoDecimals(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "roundDown",
			input:    123.456,
			expected: 123.46,
		},
		{
			name:     "roundUp",
			input:    123.455,
			expected: 123.46,
		},
		{
			name:     "noRoundingNeeded",
			input:    123.45,
			expected: 123.45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, roundToTwoDecimals(tt.input))
		})
	}
}
