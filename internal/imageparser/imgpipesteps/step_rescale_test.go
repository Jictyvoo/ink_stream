package imgpipesteps

import (
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils/testimgs"
)

func TestStepRescaleImage_PerformExec(t *testing.T) {
	imgFixtures := ([2]image.Image)(
		testimgs.ImageFixtures(2, []byte("TestStepRescaleImage_PerformExec")),
	)
	testCases := []struct {
		name                  string
		originalSize          deviceprof.Resolution
		targetSize            deviceprof.Resolution
		fillWith              color.Color
		inputImg, expectedImg image.Image
	}{
		{
			name:         "Empty input",
			originalSize: deviceprof.Resolution{Width: 0, Height: 0},
			targetSize:   deviceprof.Resolution{Width: 2, Height: 2},
		},
		{
			name: "Keep same size",
			targetSize: func() deviceprof.Resolution {
				bounds := imgFixtures[0].Bounds()
				return deviceprof.Resolution{Width: uint(bounds.Dx()), Height: uint(bounds.Dy())}
			}(),
			fillWith:    color.White,
			inputImg:    imgFixtures[0],
			expectedImg: imgFixtures[0],
		},
		{
			name:       "Make image be 0",
			inputImg:   imgFixtures[1],
			targetSize: deviceprof.Resolution{Width: 0, Height: 0},
		},
		{
			name:         "Increase image size by x6",
			originalSize: deviceprof.Resolution{Width: 7, Height: 7},
			targetSize:   deviceprof.Resolution{Width: 42, Height: 42},
			fillWith:     color.White,
		},
		{
			name:         "Multiply only width size",
			originalSize: deviceprof.Resolution{Width: 1, Height: 1},
			targetSize:   deviceprof.Resolution{Width: 6, Height: 1},
			fillWith:     color.RGBA{R: 128, G: 63, B: 16, A: 255},
		},
		{
			name:         "Divide all dimensions by 3",
			originalSize: deviceprof.Resolution{Width: 9, Height: 27},
			targetSize:   deviceprof.Resolution{Width: 3, Height: 9},
			fillWith:     color.RGBA{R: 8, G: 127, B: 31, A: 255},
		},
		{
			name:         "Target Height 8x smaller", // Must add padding before resize
			originalSize: deviceprof.Resolution{Width: 4, Height: 32},
			targetSize:   deviceprof.Resolution{Width: 4, Height: 4},
			fillWith:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			step := NewStepRescale(tCase.targetSize)
			mockImage := func(size deviceprof.Resolution, fillValue color.Color) image.Image {
				img := image.NewRGBA(
					image.Rect(0, 0, int(size.Width), int(size.Height)),
				)
				if fillValue == nil {
					fillValue = color.RGBA{}
				}
				for y := uint(0); y < size.Height; y++ {
					for x := uint(0); x < size.Width; x++ {
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

func TestCalculateNewDimensions(t *testing.T) {
	tests := []struct {
		name             string
		inputWidth       int
		inputHeight      int
		targetResolution deviceprof.Resolution
		expectedResult   [2]uint
	}{
		{
			name:             "Aspect ratio greater",
			inputWidth:       800,
			inputHeight:      600,
			targetResolution: deviceprof.Resolution{Width: 1080, Height: 720},
			expectedResult:   [2]uint{100, 0},
		},
		{
			name:             "Same Aspect Ratio",
			inputWidth:       600,
			inputHeight:      800,
			targetResolution: deviceprof.Resolution{Width: 750, Height: 1000},
			expectedResult:   [2]uint{0, 0},
		},
		{
			name:             "Zero input width and height",
			inputWidth:       0,
			inputHeight:      0,
			targetResolution: deviceprof.Resolution{Width: 1000, Height: 750},
			expectedResult:   [2]uint{1000, 750},
		},
		{
			name:             "Zero input width",
			inputWidth:       0,
			inputHeight:      800,
			targetResolution: deviceprof.Resolution{Width: 1000, Height: 750},
			expectedResult:   [2]uint{1066, 0},
		},
		{
			name:             "Invalid input width and height",
			inputWidth:       -10, // Interpreted as positive
			inputHeight:      0,
			targetResolution: deviceprof.Resolution{Width: 1000, Height: 750},
			expectedResult:   [2]uint{0, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := NewStepRescale(tt.targetResolution)
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
