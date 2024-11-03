package imageparser

import (
	"image"
	"image/color"
	"math"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils"
)

type StepGammaCorrectionImage struct{}

func NewStepGammaCorrection() StepGammaCorrectionImage {
	return StepGammaCorrectionImage{}
}

func (step StepGammaCorrectionImage) PerformExec(
	state *pipeState, opts processOptions,
) (err error) {
	bounds := state.img.Bounds()
	newImg := image.NewRGBA(bounds)

	// GammaCorrection applies gamma correction on a given value
	correction := func(a uint8) uint8 {
		return uint8(
			min(
				imgutils.MaxPixelValue,
				max(0, imgutils.MaxPixelValue*math.Pow(
					float64(a)/imgutils.MaxPixelValue, opts.gamma,
				)),
			),
		)
	}

	// GenerateLUT as a lookup table
	var lut [imgutils.MaxPixelValue + 1]uint8
	for i := 0; i < imgutils.MaxPixelValue+1; i++ {
		lut[i] = correction(uint8(i))
	}

	for x, y := range imgutils.Iterator(state.img) {
		r, g, b, a := state.img.At(x, y).RGBA()

		// Applying the LUT to each channel, assuming 8-bit image (values scaled down from 16-bit)
		newR := lut[r>>8]
		newG := lut[g>>8]
		newB := lut[b>>8]

		newImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
	}

	state.img = newImg
	return
}
