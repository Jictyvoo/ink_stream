package imgutils

import (
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
)

func TestDominantColorInRegion(t *testing.T) {
	testCases := []struct {
		name           string
		img            image.Image
		region         image.Rectangle
		optTakeAverage bool
		expected       color.NRGBA
	}{
		{
			name:     "Manga kind image",
			img:      testimgs.ImageGenericMangaPage(),
			region:   image.Rect(0, 0, 2, 106),
			expected: color.NRGBA{R: 222, G: 224, B: 223, A: 255},
		},
		{
			name:     "Solid red image",
			img:      testimgs.NewSolidImage(image.Rect(0, 0, 10, 10), color.NRGBA{R: 255, A: 255}),
			region:   image.Rect(0, 0, 10, 10),
			expected: color.NRGBA{R: 255, A: 255},
		},
		{
			name: "Two-color average",
			img: func() image.Image {
				img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
				img.Set(0, 0, color.NRGBA{R: 100, G: 100, B: 100, A: 255})
				img.Set(1, 0, color.NRGBA{R: 200, G: 200, B: 200, A: 255})
				return img
			}(),
			region:         image.Rect(0, 0, 2, 1),
			expected:       color.NRGBA{R: 150, G: 150, B: 150, A: 255},
			optTakeAverage: true,
		},
		{
			name: "With transparent pixel (ignored)",
			img: func() image.Image {
				img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
				img.Set(0, 0, color.NRGBA{R: 100, G: 150, B: 200, A: 255})
				img.Set(1, 0, color.NRGBA{}) // transparent
				return img
			}(),
			region:   image.Rect(0, 0, 2, 1),
			expected: color.NRGBA{R: 100, G: 150, B: 200, A: 255},
		},
		{
			name: "All transparent pixels",
			img: func() image.Image {
				img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
				img.Set(0, 0, color.NRGBA{R: 50, G: 50, B: 50})
				img.Set(1, 0, color.NRGBA{R: 100, G: 100, B: 100})
				return img
			}(),
			region:   image.Rect(0, 0, 2, 1),
			expected: color.NRGBA{},
		},
		{
			name: "Empty region",
			img: testimgs.NewSolidImage(
				image.Rect(0, 0, 2, 2),
				color.NRGBA{R: 10, G: 20, B: 30, A: 255},
			),
			region:   image.Rect(0, 0, 0, 0),
			expected: color.NRGBA{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := DominantColorInRegion(tt.img, tt.region, tt.optTakeAverage)
			if got != tt.expected {
				t.Errorf("DominantColorInRegion() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestImageMarginDominantColor(t *testing.T) {
	expected := Margins[color.Color]{
		Top:    color.NRGBA{R: 255, A: 255},         // Red
		Bottom: color.NRGBA{G: 255, A: 255},         // Green
		Left:   color.NRGBA{B: 255, A: 255},         // Blue
		Right:  color.NRGBA{R: 255, G: 255, A: 255}, // Yellow
	}

	img := func() image.Image {
		drawImg := image.NewNRGBA(image.Rect(0, 0, 100, 100))

		// Fill whole image with white (background)
		FillImageRegionWithColor(
			drawImg,
			image.Rect(0, 0, 100, 100),
			color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		)

		// Top 5% = red
		for x, y := range RegionIterator(image.Rect(0, 0, 100, 5)) {
			drawImg.Set(x, y, expected.Top)
		}

		// Bottom 5% = green
		for x, y := range RegionIterator(image.Rect(0, 95, 100, 100)) {
			drawImg.Set(x, y, expected.Bottom)
		}

		// Left 5% = blue
		for x, y := range RegionIterator(image.Rect(0, 5, 5, 100)) {
			drawImg.Set(x, y, expected.Left)
		}

		// Right 5% = yellow
		for x, y := range RegionIterator(image.Rect(95, 0, 100, 95)) {
			drawImg.Set(x, y, expected.Right)
		}

		return drawImg
	}()
	result := ImageMarginDominantColor(img, 10, 10, 5) // 5% margins

	// Comparison helper
	assertColorEqual := func(name string, got, want color.Color) {
		if got != want {
			t.Errorf("%s: got %v, want %v", name, got, want)
		}
	}

	assertColorEqual("Top", result.Top, expected.Top)
	assertColorEqual("Bottom", result.Bottom, expected.Bottom)
	assertColorEqual("Left", result.Left, expected.Left)
	assertColorEqual("Right", result.Right, expected.Right)
}
