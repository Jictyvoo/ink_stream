package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepAutoContrastImage)(nil)

type StepAutoContrastImage struct {
	cutoff       [2]float64
	applyColor   bool
	gammaCorrect StepGammaCorrectionImage
	imageparser.BaseImageStep
}

func NewStepAutoContrast(cutLow, cutHigh float64) *StepAutoContrastImage {
	return &StepAutoContrastImage{
		cutoff: [2]float64{cutLow, cutHigh},
	}
}

func (step StepAutoContrastImage) StepID() string {
	return "autocontrast"
}

func (step *StepAutoContrastImage) UpdateDrawFactory(fac imgutils.DrawImageFactory) {
	step.BaseImageStep.UpdateDrawFactory(fac)
	step.gammaCorrect.UpdateDrawFactory(fac)
}

// AutoContrast applies autocontrast to an image.
func (step StepAutoContrastImage) AutoContrast(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := step.DrawImage(img.ColorModel(), bounds)
	histogram := imgutils.CalculateHistogram(img)

	// Apply cutoff to histogram
	if step.cutoff != [2]float64{} {
		for i := range uint8(3) {
			newChannel := imgutils.ApplyCutoff(
				histogram.Channel(i),
				step.cutoff[0],
				step.cutoff[1],
			)
			histogram.Set(i, newChannel)
		}
	}

	var (
		// Determine minVal and maxVal values in the image
		minVal, maxVal = histogram.HiloHistogram()
		scale          [3]float64
		lookupTable    [3]imgutils.ChannelHistogram
	)

	for i := range uint8(3) {
		scale[i] = 1

		// Avoid division by zero
		if maxVal[i] != minVal[i] {
			scale[i] = imgutils.MaxPixelValue / float64(maxVal[i]-minVal[i])
		}

		// Apply auto-contrast transformation
		{
			// Offset to adjust each channel based on minVal and scale
			offset := -float64(minVal[i]) * scale[i]

			// Fill the lookup table for this channel
			for pixelIndex := 0; pixelIndex <= imgutils.MaxPixelValue; pixelIndex++ {
				// Calculate the adjusted pixel value using the scale and offset
				adjustedValue := float64(pixelIndex)*scale[i] + offset
				lookupTable[i][pixelIndex] = uint32(imgutils.NormalizePixel(adjustedValue))
			}
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

func (step StepAutoContrastImage) PerformExec(
	state *imageparser.PipeState,
	opts imageparser.ProcessOptions,
) (err error) {
	if opts.Gamma < 0.1 {
		if step.applyColor {
			opts.Gamma = 1.0
		}
	}

	if uint16(opts.Gamma) == 1 || opts.Gamma == 0 {
		state.Img = step.AutoContrast(state.Img)
		return err
	}

	if err = step.gammaCorrect.PerformExec(state, opts); err != nil {
		return err
	}
	state.Img = step.AutoContrast(state.Img)
	return err
}
