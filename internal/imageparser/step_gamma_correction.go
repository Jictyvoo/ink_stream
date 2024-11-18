package imageparser

import (
	"image/color"
	"math"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var (
	_ UnitStep = (*StepGammaCorrectionImage)(nil)
	_ PipeStep = (*StepGammaCorrectionImage)(nil)
)

type StepGammaCorrectionImage struct {
	lut   [256]uint8
	gamma float64
	baseImageStep
}

func NewStepGammaCorrection() *StepGammaCorrectionImage {
	return &StepGammaCorrectionImage{lut: StepGammaCorrectionImage{}.lookupTable(1)}
}

func NewStepGammaCorrectionPreDefined(gamma float64) StepGammaCorrectionImage {
	return StepGammaCorrectionImage{
		lut:   StepGammaCorrectionImage{}.lookupTable(gamma),
		gamma: gamma,
	}
}

func (step StepGammaCorrectionImage) lookupTable(gamma float64) [256]uint8 {
	// GammaCorrection applies gamma correction on a given value
	correction := func(a uint8) uint8 {
		return uint8(
			min(
				imgutils.MaxPixelValue,
				max(0, imgutils.MaxPixelValue*math.Pow(
					float64(a)/imgutils.MaxPixelValue, gamma,
				)),
			),
		)
	}

	// GenerateLUT as a lookup table
	var lut [imgutils.MaxPixelValue + 1]uint8
	for i := 0; i < imgutils.MaxPixelValue+1; i++ {
		lut[i] = correction(uint8(i))
	}
	return lut
}

func (step StepGammaCorrectionImage) PerformExec(
	state *pipeState, opts processOptions,
) (err error) {
	bounds := state.img.Bounds()
	newImg := step.drawImage(state.img, bounds)
	if opts.gamma != 1 {
		step.lut = step.lookupTable(opts.gamma)
	}

	for x, y := range imgutils.Iterator(state.img) {
		imgColor := state.img.At(x, y)
		newColor := step.PixelStep(imgColor)
		newImg.Set(x, y, newColor)
	}

	state.img = newImg
	return
}

func (step StepGammaCorrectionImage) PixelStep(imgColor color.Color) color.Color {
	r, g, b, a := imgColor.RGBA()

	// Applying the LUT to each channel, assuming 8-bit image (values scaled down from 16-bit)
	newR := step.lut[r>>8]
	newG := step.lut[g>>8]
	newB := step.lut[b>>8]

	return color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)}
}
