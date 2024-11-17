package imgutils

import (
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils/testimgs"
)

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
