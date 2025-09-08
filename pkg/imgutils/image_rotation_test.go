package imgutils

import (
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func TestNewOrientation(t *testing.T) {
	tests := []struct {
		name     string
		rect     image.Rectangle
		expected inktypes.ImageOrientation
	}{
		{
			name:     "Landscape orientation",
			rect:     image.Rect(0, 0, 200, 100),
			expected: inktypes.OrientationLandscape,
		},
		{
			name:     "Portrait orientation",
			rect:     image.Rect(0, 0, 100, 200),
			expected: inktypes.OrientationPortrait,
		},
		{
			name:     "Square image treated as portrait",
			rect:     image.Rect(0, 0, 100, 100),
			expected: inktypes.OrientationPortrait, // Since Dx() == Dy(), it defaults to Portrait
		},
		{
			name:     "Negative coordinates, still landscape",
			rect:     image.Rect(-100, -50, 100, 50),
			expected: inktypes.OrientationLandscape,
		},
		{
			name:     "Zero area image, portrait default",
			rect:     image.Rect(0, 0, 0, 0),
			expected: inktypes.OrientationPortrait,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrientation(tt.rect)
			if got != tt.expected {
				t.Errorf("NewOrientation(%v) = %v; want %v", tt.rect, got, tt.expected)
			}
		})
	}
}

func TestRotateImage(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		rotation       RotationDegrees
		expectedColors [][]color.Color
	}{
		{
			name:     "Rotate 90 degrees",
			rotation: Rotation90Degrees,
			expectedColors: [][]color.Color{
				{
					color.RGBA{A: 255},
					color.RGBA{R: 255, G: 255, A: 255},
					color.RGBA{R: 255, A: 255},
				},
				{
					color.RGBA{R: 255, G: 255, B: 255, A: 255},
					color.RGBA{G: 255, B: 255, A: 255},
					color.RGBA{G: 255, A: 255},
				},
				{
					color.RGBA{R: 128, G: 128, B: 128, A: 255},
					color.RGBA{R: 255, B: 255, A: 255},
					color.RGBA{B: 255, A: 255},
				},
			},
		},
		{
			name:     "Rotate 180 degrees",
			rotation: Rotation180Degrees,
			expectedColors: [][]color.Color{
				{
					color.RGBA{R: 128, G: 128, B: 128, A: 255},
					color.RGBA{R: 255, G: 255, B: 255, A: 255},
					color.RGBA{A: 255},
				},
				{
					color.RGBA{R: 255, B: 255, A: 255},
					color.RGBA{G: 255, B: 255, A: 255},
					color.RGBA{R: 255, G: 255, A: 255},
				},
				{
					color.RGBA{B: 255, A: 255},
					color.RGBA{G: 255, A: 255},
					color.RGBA{R: 255, A: 255},
				},
			},
		},
		{
			name:     "Rotate 270 degrees",
			rotation: Rotation270Degrees,
			expectedColors: [][]color.Color{
				{
					color.RGBA{B: 255, A: 255},
					color.RGBA{R: 255, B: 255, A: 255},
					color.RGBA{R: 128, G: 128, B: 128, A: 255},
				},
				{
					color.RGBA{G: 255, A: 255},
					color.RGBA{G: 255, B: 255, A: 255},
					color.RGBA{R: 255, G: 255, B: 255, A: 255},
				},
				{
					color.RGBA{R: 255, A: 255},
					color.RGBA{R: 255, G: 255, A: 255},
					color.RGBA{A: 255},
				},
			},
		},
	}

	originalImage := testimgs.ImageMultiColorSquare()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotated := RotateImage(originalImage, tt.rotation)

			// Validate the rotated image
			for y := 0; y < len(tt.expectedColors); y++ {
				for x := 0; x < len(tt.expectedColors[y]); x++ {
					if got, want := rotated.At(x, y), tt.expectedColors[y][x]; got != want {
						t.Errorf("At (%d, %d): got %v, want %v", x, y, got, want)
					}
				}
			}
		})
	}
}
