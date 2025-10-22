package inktypes

import "testing"

func TestNewOrientation(t *testing.T) {
	tests := []struct {
		name     string
		dim      ImageDimensions
		expected ImageOrientation
	}{
		{
			name:     "portrait",
			dim:      ImageDimensions{Width: 800, Height: 1200},
			expected: OrientationPortrait,
		},
		{
			name:     "landscape",
			dim:      ImageDimensions{Width: 1200, Height: 800},
			expected: OrientationLandscape,
		},
		{
			name:     "square",
			dim:      ImageDimensions{Width: 1000, Height: 1000},
			expected: OrientationPortrait,
		},
		{
			name:     "zero width",
			dim:      ImageDimensions{Width: 0, Height: 1000},
			expected: OrientationPortrait,
		},
		{
			name:     "zero height",
			dim:      ImageDimensions{Width: 1000, Height: 0},
			expected: OrientationLandscape,
		},
		{
			name:     "both zero",
			dim:      ImageDimensions{Width: 0, Height: 0},
			expected: OrientationPortrait,
		},
		{
			name:     "maximum values",
			dim:      ImageDimensions{Width: 65535, Height: 65535},
			expected: OrientationPortrait,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewOrientation(tt.dim)
			if result != tt.expected {
				t.Errorf("NewOrientation(%v) = %v, want %v", tt.dim, result, tt.expected)
			}
		})
	}
}
