package formatter

import "testing"

func TestFormatUSD(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "normal use",
			input:       "1234567.81231231243",
			expected:    "$1,234,567.81",
			expectError: false,
		},
		{
			name:        "zero",
			input:       "0",
			expected:    "$0.00",
			expectError: false,
		},
		{
			name:        "negative",
			input:       "-1200.56",
			expected:    "$-1,200.56",
			expectError: false,
		},
		{
			name:        "invalid input",
			input:       "invalid",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatUSD(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("FormatUSD: for %v expected error = %v, got = %v", tt.input, tt.expectError, err != nil)
				return
			}
			if result != tt.expected {
				t.Errorf("FormatUSD: for %v expected = %v, got = %v", tt.input, tt.expected, result)
			}
		})
	}
}
