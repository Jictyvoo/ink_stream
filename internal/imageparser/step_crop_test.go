package imageparser

import (
	_ "embed"
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
)

func TestStepCropRotateImage_PerformExec(t *testing.T) {
	tests := []struct {
		name           string
		initialBounds  image.Rectangle
		rotateImage    bool
		expectedBounds image.Rectangle
	}{
		{
			name:           "Rotate Portrait Image",
			initialBounds:  image.Rect(0, 0, 200, 100),
			rotateImage:    true,
			expectedBounds: image.Rect(0, 0, 100, 200), // Dimensions swapped after rotation
		},
		{
			name:           "Crop Portrait Image without Rotation",
			initialBounds:  image.Rect(0, 0, 200, 100),
			rotateImage:    false,
			expectedBounds: image.Rect(100, 0, 200, 100), // Cropped to square dimensions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a solid-colored test image
			img := testimgs.NewSolidImage(tt.initialBounds, color.White)
			state := &pipeState{img: img}

			// Instantiate the StepCropRotateImage step
			step := NewStepCropRotate(
				tt.rotateImage,
				color.Palette{color.Black, color.White},
				imgutils.OrientationPortrait,
			)

			// Execute the step
			if err := step.PerformExec(state, processOptions{}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate the resulting image bounds
			resultBounds := state.img.Bounds()
			if resultBounds != tt.expectedBounds {
				t.Errorf(
					"unexpected bounds: got %v, want %v",
					resultBounds, tt.expectedBounds,
				)
			}
		})
	}
}
