package tools

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "valid input with comma",
			input:    "123,45",
			expected: 123.45,
		},
		{
			name:     "valid input with dot",
			input:    "678.91",
			expected: 678.91,
		},
		{
			name:     "valid input with no decimal",
			input:    "42",
			expected: 42.0,
		},
		{
			name:     "invalid input with letters",
			input:    "abc",
			expected: 0.0,
		},
		{
			name:     "empty input",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "input with multiple commas",
			input:    "1,234,567.89",
			expected: 0.0,
		},
		{
			name:     "input with mixed separators",
			input:    "1,234.567",
			expected: 1234.567,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("ParseFloat(%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
