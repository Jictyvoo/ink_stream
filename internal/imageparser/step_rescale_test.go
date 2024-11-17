package imageparser

import (
	"image"
	"image/color"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/deviceprof"
	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils"
	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils/testimgs"
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
			name:         "Divide Height by 8",
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
				state = pipeState{img: tCase.inputImg}
				opts  processOptions
			)

			if err := step.PerformExec(&state, opts); err != nil {
				t.Fatalf("%s: PerformExec: %v", tCase.name, err.Error())
			}

			result := state.img
			if !imgutils.IsImageEqual(result, tCase.expectedImg) {
				t.Errorf(
					"expected: %#v, actual: %#v", tCase.expectedImg, result,
				)
			}
		})
	}
}
