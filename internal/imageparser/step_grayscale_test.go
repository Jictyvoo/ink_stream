package imageparser

import (
	"image"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgcompare"
)

func TestStepGrayScaleImage_PerformExec(t *testing.T) {
	// Define the custom color image and its expected grayscale counterpart
	customColorImage := &image.RGBA{
		Pix: []uint8{
			0xff, 0x0, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0x0, 0xff, 0xff,
			0xff, 0xff, 0x0, 0xff, 0x0, 0xff, 0xff, 0xff, 0xff, 0x0, 0xff, 0xff,
			0x80, 0x80, 0x80, 0xff, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff,
		},
		Stride: 12,
		Rect:   image.Rect(0, 0, 3, 3),
	}
	expectedGrayImage := &image.Gray{
		Pix:    []uint8{0x4c, 0x96, 0x1d, 0xe2, 0xb3, 0x69, 0x80, 0x0, 0xff},
		Stride: 3,
		Rect:   image.Rect(0, 0, 3, 3),
	}

	testCases := []struct {
		name        string
		inputImg    image.Image
		expectedImg image.Image
	}{
		{
			name:        "Single pixel",
			inputImg:    image.NewRGBA(image.Rect(0, 0, 1, 1)),
			expectedImg: image.NewGray(image.Rect(0, 0, 1, 1)),
		},
		{
			name:        "3x3 Custom Image",
			inputImg:    customColorImage,
			expectedImg: expectedGrayImage,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			step := NewStepGrayScale()
			var (
				state = pipeState{img: tCase.inputImg}
				opts  processOptions
			)

			// Perform grayscale conversion
			if err := step.PerformExec(&state, opts); err != nil {
				t.Fatalf("PerformExec: %v", err.Error())
			}

			// Validate that the output matches the expected grayscale image
			result := state.img
			if !imgcompare.IsImageEqual(result, tCase.expectedImg) {
				t.Errorf(
					"expected: %#v, actual: %#v", tCase.expectedImg, result,
				)
			}
		})
	}
}
