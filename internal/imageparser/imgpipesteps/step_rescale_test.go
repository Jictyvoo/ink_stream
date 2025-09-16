package imgpipesteps

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func TestStepRescaleImage_PerformExec(t *testing.T) {
	imgFixtures := ([2]image.Image)(
		testimgs.ImageFixtures(2, []byte("TestStepRescaleImage_PerformExec")),
	)
	testCases := []struct {
		name                  string
		originalSize          inktypes.ImageDimensions
		targetSize            inktypes.ImageDimensions
		fillWith              color.Color
		inputImg, expectedImg image.Image
	}{
		{
			name:         "Empty input",
			originalSize: inktypes.ImageDimensions{Width: 0, Height: 0},
			targetSize:   inktypes.ImageDimensions{Width: 2, Height: 2},
		},
		{
			name: "Keep same size",
			targetSize: func() inktypes.ImageDimensions {
				bounds := imgFixtures[0].Bounds()
				return inktypes.ImageDimensions{
					Width:  uint16(bounds.Dx()),
					Height: uint16(bounds.Dy()),
				}
			}(),
			fillWith:    color.White,
			inputImg:    imgFixtures[0],
			expectedImg: imgFixtures[0],
		},
		{
			name:       "Make image be 0",
			inputImg:   imgFixtures[1],
			targetSize: inktypes.ImageDimensions{Width: 0, Height: 0},
		},
		{
			name:         "Increase image size by x6",
			originalSize: inktypes.ImageDimensions{Width: 7, Height: 7},
			targetSize:   inktypes.ImageDimensions{Width: 42, Height: 42},
			fillWith:     color.White,
		},
		{
			name:         "Multiply only width size",
			originalSize: inktypes.ImageDimensions{Width: 1, Height: 1},
			targetSize:   inktypes.ImageDimensions{Width: 6, Height: 1},
			fillWith:     color.RGBA{R: 128, G: 63, B: 16, A: 255},
		},
		{
			name:         "Divide all dimensions by 3",
			originalSize: inktypes.ImageDimensions{Width: 9, Height: 27},
			targetSize:   inktypes.ImageDimensions{Width: 3, Height: 9},
			fillWith:     color.RGBA{R: 8, G: 127, B: 31, A: 255},
		},
		{
			name:         "Target Height 8x smaller", // Must add padding before resize
			originalSize: inktypes.ImageDimensions{Width: 4, Height: 32},
			targetSize:   inktypes.ImageDimensions{Width: 4, Height: 4},
			fillWith:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			step := NewStepRescale(tCase.targetSize, true)
			mockImage := func(size inktypes.ImageDimensions, fillValue color.Color) image.Image {
				img := image.NewRGBA(
					image.Rect(0, 0, int(size.Width), int(size.Height)),
				)
				if fillValue == nil {
					fillValue = color.RGBA{}
				}
				for y := uint16(0); y < size.Height; y++ {
					for x := uint16(0); x < size.Width; x++ {
						img.Set(int(x), int(y), fillValue)
					}
				}

				return img
			}

			// Generate the original image with specified size, filled with opaque white pixels
			if tCase.inputImg == nil {
				tCase.inputImg = mockImage(tCase.originalSize, tCase.fillWith)
			}

			// Generate the expected output image with the target size, also filled with opaque white pixels
			if tCase.expectedImg == nil {
				tCase.expectedImg = mockImage(tCase.targetSize, tCase.fillWith)
			}
			var (
				state = imageparser.PipeState{Img: tCase.inputImg}
				opts  imageparser.ProcessOptions
			)

			if err := step.PerformExec(&state, opts); err != nil {
				t.Fatalf("%s: PerformExec: %v", tCase.name, err.Error())
			}

			result := state.Img
			if !imgutils.IsImageEqual(result, tCase.expectedImg) {
				t.Errorf(
					"expected: %#v, actual: %#v", tCase.expectedImg, result,
				)
			}
		})
	}
}

func TestStepRescaleImage_updateTargetResolution(t *testing.T) {
	testCases := []struct {
		name     string
		original inktypes.ImageDimensions
		target   inktypes.ImageDimensions
		expected inktypes.ImageDimensions
	}{
		{
			name:     "Landscape to Portrait - Simple scaling",
			original: inktypes.ImageDimensions{Width: 651, Height: 1008},
			target:   inktypes.ImageDimensions{Width: 1072, Height: 1448},
			expected: inktypes.ImageDimensions{Width: 935, Height: 1448},
		},
		{
			name:     "Portrait to Square - Original narrower than target",
			original: inktypes.ImageDimensions{Width: 100, Height: 200},
			target:   inktypes.ImageDimensions{Width: 50, Height: 50},
			expected: inktypes.ImageDimensions{Width: 25, Height: 50},
		},
		{
			name:     "Landscape to Landscape - Original wider than target",
			original: inktypes.ImageDimensions{Width: 300, Height: 200},
			target:   inktypes.ImageDimensions{Width: 200, Height: 100},
			expected: inktypes.ImageDimensions{Width: 150, Height: 100},
		},
		{
			name:     "Portrait to Landscape - Original taller than target",
			original: inktypes.ImageDimensions{Width: 200, Height: 300},
			target:   inktypes.ImageDimensions{Width: 200, Height: 100},
			expected: inktypes.ImageDimensions{Width: 66, Height: 100},
		},
		{
			name:     "Square to Square - Upscaling",
			original: inktypes.ImageDimensions{Width: 200, Height: 200},
			target:   inktypes.ImageDimensions{Width: 250, Height: 250},
			expected: inktypes.ImageDimensions{Width: 250, Height: 250},
		},
		{
			name:     "Portrait with zero width",
			original: inktypes.ImageDimensions{Width: 0, Height: 200},
			target:   inktypes.ImageDimensions{Width: 100, Height: 100},
			expected: inktypes.ImageDimensions{Width: 0, Height: 100},
		},
		{
			name:     "Landscape with zero height",
			original: inktypes.ImageDimensions{Width: 200, Height: 0},
			target:   inktypes.ImageDimensions{Width: 100, Height: 100},
			expected: inktypes.ImageDimensions{Width: 100, Height: 0},
		},
		{
			name:     "Any to Zero target dimensions",
			original: inktypes.ImageDimensions{Width: 200, Height: 300},
			target:   inktypes.ImageDimensions{Width: 0, Height: 0},
			expected: inktypes.ImageDimensions{Width: 0, Height: 0},
		},
		{
			name:     "Portrait to Landscape - Scaling down",
			original: inktypes.ImageDimensions{Width: 200, Height: 300},
			target:   inktypes.ImageDimensions{Width: 400, Height: 200},
			expected: inktypes.ImageDimensions{Width: 133, Height: 200},
		},
	}

	for _, tt := range testCases {
		if tt.name == "" {
			tt.name = fmt.Sprintf("%+v->%+v", tt.original, tt.target)
		}
		t.Run(tt.name, func(t *testing.T) {
			step := StepRescaleImage{resolution: tt.target}
			result := step.updateTargetResolution(tt.original)

			if result.Width != tt.expected.Width || result.Height != tt.expected.Height {
				t.Errorf(
					"For original dimensions (%+v) and target resolution (%+v)\n\tExpected new dimensions (%+v)\n\tGot (%+v)",
					tt.original,
					tt.target,
					tt.expected,
					result,
				)
			}
		})
	}
}
