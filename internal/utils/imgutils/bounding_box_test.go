package imgutils

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
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
			img: func(size, margin int) image.Image {
				img := image.NewRGBA(image.Rect(0, 0, size, size))
				// Fill the entire image with white
				draw.Draw(
					img, img.Bounds(),
					&image.Uniform{C: color.White},
					image.Point{}, draw.Src,
				)
				// Draw the black square in the center
				blackRect := image.Rect(margin, margin, size-margin, size-margin)
				draw.Draw(img, blackRect, &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
				return img
			}(50, 10),
			opts: BoxEliminateMinimumColor,
			// Expect bounding box to exclude the white margins
			expected: image.Rect(10, 10, 40, 40),
		},
		{
			name: "Circle with transparent background",
			img: func(size int) image.Image {
				img := image.NewRGBA(image.Rect(0, 0, size, size))
				// Fill the image with transparent pixels
				draw.Draw(
					img, img.Bounds(),
					&image.Uniform{C: color.Transparent},
					image.Point{}, draw.Src,
				)
				// Draw a filled circle in the center
				center := size / 2
				radius := center - 10
				for y := 0; y < size; y++ {
					for x := 0; x < size; x++ {
						dx, dy := x-center, y-center
						if dx*dx+dy*dy <= radius*radius {
							img.Set(x, y, color.Black)
						}
					}
				}
				return img
			}(100),
			opts: BoxEliminateTransparent,
			// Expect bounding box to fit the circle tightly
			expected: image.Rect(10, 10, 91, 91),
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
	cropped := image.NewRGBA(rect)
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)
	return cropped
}
