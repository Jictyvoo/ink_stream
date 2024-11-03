package imageparser

import (
	"image"
	"image/color"
	"math"
)

type StepAutoContrastImage struct{}

func NewStepAutoContrast() StepAutoContrastImage {
	return StepAutoContrastImage{}
}

// GammaCorrect applies gamma correction to an image.
func (sgsi StepAutoContrastImage) GammaCorrect(img image.Image, gamma float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// GammaCorrection applies gamma correction on a given value
	correction := func(a uint8) uint8 {
		return uint8(
			math.Min(
				maxPixelValue,
				math.Max(0, maxPixelValue*math.Pow(float64(a)/maxPixelValue, gamma)),
			),
		)
	}

	// GenerateLUT as a lookup table
	var lut [maxPixelValue + 1]uint8
	for i := 0; i < maxPixelValue+1; i++ {
		lut[i] = correction(uint8(i))
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Applying the LUT to each channel, assuming 8-bit image (values scaled down from 16-bit)
			newR := lut[r>>8]
			newG := lut[g>>8]
			newB := lut[b>>8]

			newImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return newImg
}

// AutoContrast applies autocontrast to an image.
func (sgsi StepAutoContrastImage) AutoContrast(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	var (
		minVal, maxVal = [3]uint8{}, [3]uint8{maxPixelValue, maxPixelValue, maxPixelValue}
		stopChannels   [3]struct{ min, max bool }
		histogram      = calculateHistogram(img)
	)

	// Determine minVal and maxVal values in the image
	minVal, maxVal = histogram.hiloHistogram(minVal, maxVal, stopChannels)

	// Avoid division by zero
	for i := range len(minVal) {
		if maxVal[i] == minVal[i] {
			maxVal[i] = maxPixelValue
			maxVal[i] = 0
		}
	}

	// Apply autocontrast transformation
	scale := [3]float64{
		maxPixelValue / float64(maxVal[0]-minVal[0]),
		maxPixelValue / float64(maxVal[1]-minVal[1]),
		maxPixelValue / float64(maxVal[2]-minVal[2]),
	}

	clamp := func(value float64) uint8 {
		if value < 0 {
			return 0
		} else if value > 255 {
			return 255
		}
		return uint8(value)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Adjust each channel individually, using `clamp` to keep values within the 0-255 range.
			newR := clamp(scale[0] * float64(r>>8-uint32(minVal[0])))
			newG := clamp(scale[1] * float64(g>>8-uint32(minVal[1])))
			newB := clamp(scale[2] * float64(b>>8-uint32(minVal[2])))

			newImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a / 257)})
		}
	}

	return newImg
}

func (sgsi StepAutoContrastImage) PerformExec(state *pipeState, opts processOptions) (err error) {
	if opts.gamma < 0.1 {
		if opts.applyColor {
			opts.gamma = 1.0
		}
	}

	if uint16(opts.gamma) == 1 || opts.gamma == 0 {
		state.img = sgsi.AutoContrast(state.img)
		return
	}
	state.img = sgsi.AutoContrast(sgsi.GammaCorrect(state.img, opts.gamma))
	return
}
