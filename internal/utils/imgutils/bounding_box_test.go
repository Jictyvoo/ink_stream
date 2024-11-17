package imgutils

import (
	"image"
	"image/draw"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils/testimgs"
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

// Helper function: Crop an image to the given rectangle
func cropImage(img image.Image, rect image.Rectangle) image.Image {
	cropped := NewDrawFromImgColorModel(img, rect)
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)
	return cropped
}
