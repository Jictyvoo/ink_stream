package inktypes

import "testing"

func TestNewReadDirection(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected ReadDirection
	}{
		{
			name:     "rtl lowercase",
			value:    "rtl",
			expected: ReadRightToLeft,
		},
		{
			name:     "rtl uppercase",
			value:    "RTL",
			expected: ReadRightToLeft,
		},
		{
			name:     "rtl mixed case",
			value:    "RtL",
			expected: ReadRightToLeft,
		},
		{
			name:     "ltr lowercase",
			value:    "ltr",
			expected: ReadLeftToRight,
		},
		{
			name:     "ltr uppercase",
			value:    "LTR",
			expected: ReadLeftToRight,
		},
		{
			name:     "ltr mixed case",
			value:    "LtR",
			expected: ReadLeftToRight,
		},
		{
			name:     "unknown value",
			value:    "unknown",
			expected: ReadUnknown,
		},
		{
			name:     "empty string",
			value:    "",
			expected: ReadUnknown,
		},
		{
			name:     "whitespace string",
			value:    "   ",
			expected: ReadUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewReadDirection(tt.value)
			if result != tt.expected {
				t.Errorf("NewReadDirection(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestReadDirectionString(t *testing.T) {
	tests := []struct {
		name      string
		direction ReadDirection
		expected  string
	}{
		{
			name:      "ReadRightToLeft",
			direction: ReadRightToLeft,
			expected:  "rtl",
		},
		{
			name:      "ReadLeftToRight",
			direction: ReadLeftToRight,
			expected:  "ltr",
		},
		{
			name:      "ReadUnknown",
			direction: ReadUnknown,
			expected:  "ltr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.direction.String()
			if result != tt.expected {
				t.Errorf("ReadDirection.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestReadDirectionStringConsistency(t *testing.T) {
	// Test that the string representation is consistent with the constructor
	tests := []struct {
		name      string
		direction ReadDirection
	}{
		{
			name:      "ReadRightToLeft",
			direction: ReadRightToLeft,
		},
		{
			name:      "ReadLeftToRight",
			direction: ReadLeftToRight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalString := tt.direction.String()
			reconstructed := NewReadDirection(originalString)
			if reconstructed != tt.direction {
				t.Errorf("Round-trip failed: %v.String() = %q, NewReadDirection(%q) = %v",
					tt.direction, originalString, originalString, reconstructed)
			}
		})
	}
}
