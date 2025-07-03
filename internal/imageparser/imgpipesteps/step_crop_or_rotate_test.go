package imgpipesteps

import (
	_ "embed"
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
)

func TestStepCropOrRotateImage_PerformExec(t *testing.T) {
	tests := []struct {
		name                string
		initialBounds       image.Rectangle
		rotateImage         bool
		expectedOrientation imgutils.ImageOrientation
		expectedBounds      []image.Rectangle
	}{
		{
			name:                "Rotate Portrait Image",
			initialBounds:       image.Rect(0, 0, 200, 100),
			rotateImage:         true,
			expectedOrientation: imgutils.OrientationPortrait,
			expectedBounds: []image.Rectangle{
				image.Rect(0, 0, 100, 200),
			}, // Dimensions swapped after rotation
		},
		{
			name:                "Crop Portrait Image without Rotation",
			initialBounds:       image.Rect(0, 0, 200, 100),
			rotateImage:         false,
			expectedOrientation: imgutils.OrientationPortrait,
			expectedBounds: []image.Rectangle{
				image.Rect(0, 0, 100, 100),
				image.Rect(100, 0, 200, 100),
			}, // Cropped to square dimensions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a solid-colored test image
			img := testimgs.NewSolidImage(tt.initialBounds, color.White)
			state := &imageparser.PipeState{Img: img}

			// Instantiate the StepCropOrRotateImage step
			step := NewStepCropOrRotate(
				tt.rotateImage, color.Palette{color.Black, color.White}, tt.expectedOrientation,
			)

			// Execute the step
			if err := step.PerformExec(state, imageparser.ProcessOptions{}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate the resulting image bounds
			resultBounds := state.Img.Bounds()
			compareImgs := state.SubImages
			if len(tt.expectedBounds) < 2 {
				compareImgs = append(compareImgs, state.Img)
			}

			for index, expected := range tt.expectedBounds {
				if compareImgs[index].Bounds() != expected {
					t.Errorf(
						"unexpected bounds: got %v, want %v",
						resultBounds, tt.expectedBounds,
					)
				}
			}
		})
	}
}
