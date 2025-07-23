package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		isValid bool
	}{
		{
			name:    "validNumber",
			number:  "4561261212345467",
			isValid: true,
		},
		{
			name:    "invalidNumber",
			number:  "4561261212345464",
			isValid: false,
		},
		{
			name:    "emptyString",
			number:  "",
			isValid: false,
		},
		{
			name:    "nonNumeric",
			number:  "abc",
			isValid: false,
		},
		{
			name:    "invalidCharacters",
			number:  "4561-2612-1234-5467",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, isValidLuhn(tt.number))
		})
	}
}

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
