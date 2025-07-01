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
