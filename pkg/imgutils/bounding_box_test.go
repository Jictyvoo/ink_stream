package imgutils

import (
	"image"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

// TestCropBox tests the CropBox function
func TestCropBox(t *testing.T) {
	tests := []struct {
		name     string
		img      image.Image
		opts     BoxOptions
		expected image.Rectangle
	}{
		{
			name: "Black square with white margins",
			img:  testimgs.ImageBlackSquareWhiteMargin(),
			opts: BoxEliminateMinimumColor,
			// Expect bounding box to exclude the white margins
			expected: image.Rect(10, 10, 40, 40),
		},
		{
			name: "Circle with transparent background",
			img:  testimgs.ImageBlackCircleWithTransparentBackground(),
			opts: BoxEliminateTransparent,
			// Expect bounding box to fit the circle tightly
			expected: image.Rect(10, 10, 91, 91),
		},
		{
			name: "Square with white on left, green, and gray margins",
			img:  testimgs.ImageBlackSquareGreenRight(false),
			opts: BoxEliminateMinimumColor,
			// Expect bounding box to exclude most of the left white margin,
			// keep only a small part of the right white margin after the green region.
			expected: image.Rect(10, 0, 45, 50),
		},
		{
			name: "Square with white on left, green, and white margin on right",
			img:  testimgs.ImageBlackSquareGreenRight(true),
			opts: BoxEliminateMinimumColor,
			// Expect bounding box to exclude most of the left white margin,
			// keep only a small part of the right white margin after the green region.
			expected: image.Rect(0, 0, 45, 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run CropBox
			actual := CropBox(tt.img, nil, tt.opts)

			// Validate bounding box
			if actual != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}

func TestMarginBox(t *testing.T) {
	tests := []struct {
		name     string
		bounds   image.Rectangle
		expected image.Rectangle
	}{
		{
			name:     "Small rectangle",
			bounds:   image.Rect(10, 10, 50, 50),
			expected: image.Rect(8, 8, 52, 52), // 5% margin added
		},
		{
			name:     "Large rectangle",
			bounds:   image.Rect(100, 200, 600, 800),
			expected: image.Rect(75, 175, 625, 825), // 5% margin added
		},
		{
			name:     "Zero size rectangle",
			bounds:   image.Rect(0, 0, 0, 0),
			expected: image.Rect(0, 0, 0, 0), // No margin as bounds are zero
		},
		{
			name:     "Rectangle near origin",
			bounds:   image.Rect(1, 1, 50, 50),
			expected: image.Rect(0, 0, 52, 52), // Margin clipped at 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run MarginBox
			actual := MarginBox(tt.bounds, 0.05)

			// Validate the bounding box
			if actual != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}

func TestHalfSplit(t *testing.T) {
	tests := []struct {
		name        string
		inputRect   image.Rectangle
		orientation inktypes.ImageOrientation
		expected    Margins[image.Rectangle]
	}{
		{
			name:        "Landscape split",
			inputRect:   image.Rect(0, 0, 100, 50),
			orientation: inktypes.OrientationLandscape,
			expected: Margins[image.Rectangle]{
				Left:  image.Rect(0, 0, 50, 50),
				Right: image.Rect(50, 0, 100, 50),
			},
		},
		{
			name:        "Portrait split",
			inputRect:   image.Rect(0, 0, 50, 100),
			orientation: inktypes.OrientationPortrait,
			expected: Margins[image.Rectangle]{
				Top:    image.Rect(0, 0, 50, 50),
				Bottom: image.Rect(0, 50, 50, 100),
			},
		},
		{
			name:        "Odd width landscape split",
			inputRect:   image.Rect(0, 0, 101, 50),
			orientation: inktypes.OrientationLandscape,
			expected: Margins[image.Rectangle]{
				Left:  image.Rect(0, 0, 50, 50),
				Right: image.Rect(50, 0, 101, 50),
			},
		},
		{
			name:        "Odd height portrait split",
			inputRect:   image.Rect(0, 0, 50, 101),
			orientation: inktypes.OrientationPortrait,
			expected: Margins[image.Rectangle]{
				Top:    image.Rect(0, 0, 50, 50),
				Bottom: image.Rect(0, 50, 50, 101),
			},
		},
		{
			name:        "Zero area rectangle",
			inputRect:   image.Rect(0, 0, 0, 0),
			orientation: inktypes.OrientationLandscape,
			expected: Margins[image.Rectangle]{
				Left:  image.Rect(0, 0, 0, 0),
				Right: image.Rect(0, 0, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HalfSplit(tt.inputRect, tt.orientation)
			if got != tt.expected {
				t.Errorf("HalfSplit() = %+v, expected %+v", got, tt.expected)
			}
		})
	}
}
