package imgpipesteps

import (
	"image"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func TestCalculateNewDimensions(t *testing.T) {
	tests := []struct {
		name             string
		inputWidth       int
		inputHeight      int
		targetResolution inktypes.ImageDimensions
		expectedResult   [2]uint
	}{
		{
			name:             "Aspect ratio greater",
			inputWidth:       800,
			inputHeight:      600,
			targetResolution: inktypes.ImageDimensions{Width: 1080, Height: 720},
			expectedResult:   [2]uint{100, 0},
		},
		{
			name:             "Same Aspect Ratio",
			inputWidth:       600,
			inputHeight:      800,
			targetResolution: inktypes.ImageDimensions{Width: 750, Height: 1000},
			expectedResult:   [2]uint{0, 0},
		},
		{
			name:             "Zero input width and height",
			inputWidth:       0,
			inputHeight:      0,
			targetResolution: inktypes.ImageDimensions{Width: 1000, Height: 750},
			expectedResult:   [2]uint{1000, 750},
		},
		{
			name:             "Zero input width",
			inputWidth:       0,
			inputHeight:      800,
			targetResolution: inktypes.ImageDimensions{Width: 1000, Height: 750},
			expectedResult:   [2]uint{1066, 0},
		},
		{
			name:             "Invalid input width and height",
			inputWidth:       -10, // Interpreted as positive
			inputHeight:      0,
			targetResolution: inktypes.ImageDimensions{Width: 1000, Height: 750},
			expectedResult:   [2]uint{0, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := NewStepMarginWrap(tt.targetResolution)
			result := step.calculateNewDimensions(
				image.Rect(0, 0, tt.inputWidth, tt.inputHeight),
			)
			if result.w != tt.expectedResult[0] || result.h != tt.expectedResult[1] {
				t.Errorf(
					"Expected [w:%d h:%d], but got [w:%d h:%d]",
					tt.expectedResult[0],
					tt.expectedResult[1],
					result.w,
					result.h,
				)
			}
		})
	}
}
