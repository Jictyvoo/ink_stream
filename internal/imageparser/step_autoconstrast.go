package imageparser

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils"
)

type StepAutoContrastImage struct {
	cutoff       [2]float64
	gammaCorrect StepGammaCorrectionImage
}

func NewStepAutoContrast(cutLow, cutHigh float64) StepAutoContrastImage {
	return StepAutoContrastImage{
		cutoff: [2]float64{cutLow, cutHigh},
	}
}

// AutoContrast applies autocontrast to an image.
func (step StepAutoContrastImage) AutoContrast(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := createDrawImage(img, bounds)

	var (
		minVal, maxVal = [3]uint8{}, [3]uint8{
			imgutils.MaxPixelValue,
			imgutils.MaxPixelValue,
			imgutils.MaxPixelValue,
		}
		histogram = imgutils.CalculateHistogram(img)
	)

	// Apply cutoff to histogram
	if step.cutoff != [2]float64{} {
		for i := range uint8(3) {
			newChannel := imgutils.ApplyCutoff(histogram.Channel(i), step.cutoff[0], step.cutoff[1])
			histogram.Set(i, newChannel)
		}
	}
	// Determine minVal and maxVal values in the image
	minVal, maxVal = histogram.HiloHistogram(minVal, maxVal)

	// Avoid division by zero
	for i := range uint8(3) {
		if maxVal[i] == minVal[i] {
			maxVal[i] = imgutils.MaxPixelValue
			maxVal[i] = 0
		}
	}

	// Apply auto-contrast transformation
	scale := [3]float64{
		imgutils.MaxPixelValue / float64(maxVal[0]-minVal[0]),
		imgutils.MaxPixelValue / float64(maxVal[1]-minVal[1]),
		imgutils.MaxPixelValue / float64(maxVal[2]-minVal[2]),
	}

	clamp := func(value float64) uint8 {
		if value < 0 {
			return 0
		} else if value > imgutils.MaxPixelValue {
			return imgutils.MaxPixelValue
		}
		return uint8(value)
	}

	var lookupTable [3]imgutils.ChannelHistogram
	for index := range 3 {
		// Offset to adjust each channel based on minVal and scale
		offset := -float64(minVal[index]) * scale[index]

		// Fill the lookup table for this channel
		for pixelIndex := 0; pixelIndex <= imgutils.MaxPixelValue; pixelIndex++ {
			// Calculate the adjusted pixel value using the scale and offset
			adjustedValue := float64(pixelIndex)*scale[index] + offset
			lookupTable[index][pixelIndex] = uint32(clamp(adjustedValue))
		}
	}
	for x, y := range imgutils.Iterator(img) {
		r, g, b, a := img.At(x, y).RGBA()
		// Adjust each channel individually, using `clamp` to keep values within the 0-255 range.
		var (
			rgb    = [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
			newRGB [3]uint8
		)
		for index := range 3 {
			// Map the original pixel value to the new value using the lookup table
			newRGB[index] = uint8(lookupTable[index][rgb[index]])
		}

		newImg.Set(x, y, color.RGBA{R: newRGB[0], G: newRGB[1], B: newRGB[2], A: uint8(a >> 8)})
	}

	return newImg
}

func (step StepAutoContrastImage) PerformExec(state *pipeState, opts processOptions) (err error) {
	if opts.gamma < 0.1 {
		if opts.applyColor {
			opts.gamma = 1.0
		}
	}

	if uint16(opts.gamma) == 1 || opts.gamma == 0 {
		state.img = step.AutoContrast(state.img)
		return
	}

	if err = step.gammaCorrect.PerformExec(state, opts); err != nil {
		return
	}
	state.img = step.AutoContrast(state.img)
	return
}
