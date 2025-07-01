package imgpipesteps

import (
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
)

func TestStepAutoCropImage_PerformExec(t *testing.T) {
	tests := []struct {
		name           string
		imgSize        image.Rectangle
		margins        imgutils.Margins[int]
		expectedBounds image.Rectangle
	}{
		{
			name:    "Crop solid border with margin",
			imgSize: image.Rect(0, 0, 100, 100),
			margins: imgutils.Margins[int]{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			},
			expectedBounds: image.Rect(9, 9, 91, 91), // cropped + wrapped
		},
		{
			name:           "No crop needed (content fills image)",
			imgSize:        image.Rect(0, 0, 100, 100),
			expectedBounds: image.Rect(0, 0, 100, 100),
		},
		{
			name:    "Avoid over-cropping if content too small",
			imgSize: image.Rect(0, 0, 100, 100),
			margins: imgutils.Margins[int]{
				Top:    10,
				Bottom: 40,
				Left:   10,
				Right:  40,
			}, // content area too small
			expectedBounds: image.Rect(0, 0, 100, 100),
		},
	}

	palette := color.Palette{
		color.White, // background
		color.Black, // foreground
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the test image with explicit margins
			img := testimgs.NewBorderedImage(
				tt.imgSize,
				tt.margins.Top, tt.margins.Bottom, tt.margins.Left, tt.margins.Right,
				color.White, color.Black,
			)

			state := &imageparser.PipeState{Img: img}
			step := NewStepAutoCrop(palette)

			if err := step.PerformExec(state, imageparser.ProcessOptions{}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotBounds := state.Img.Bounds()
			if gotBounds != tt.expectedBounds {
				t.Errorf("unexpected result bounds: got %v, want %v", gotBounds, tt.expectedBounds)
			}
		})
	}
}
