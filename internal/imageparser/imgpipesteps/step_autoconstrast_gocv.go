//go:build cgo

package imgpipesteps

import (
	"gocv.io/x/gocv"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepAutoContrastImage)(nil)

type StepAutoContrastImage struct {
	cutoff       [2]float64
	gammaCorrect StepGammaCorrectionImage
	imageparser.BaseImageStep
}

func NewStepAutoContrast(cutLow, cutHigh float64) *StepAutoContrastImage {
	return &StepAutoContrastImage{
		cutoff: [2]float64{cutLow, cutHigh},
	}
}

func (step *StepAutoContrastImage) UpdateDrawFactory(fac imgutils.DrawImageFactory) {
	step.BaseImageStep.UpdateDrawFactory(fac)
	step.gammaCorrect.UpdateDrawFactory(fac)
}

// AutoContrast applies the Stylization function for auto-contrast-like effect.
func (step StepAutoContrastImage) AutoContrast(img gocv.Mat) (gocv.Mat, error) {
	// Create a destination Mat to hold the result
	dst := gocv.NewMat()
	defer func() {
		if dst.Empty() {
			dst.Close()
		}
	}()

	// Apply the Stylization function
	gocv.Stylization(img, &dst, 10, 0.45)
	return dst, nil
}

func (step StepAutoContrastImage) PerformExec(
	state *imageparser.PipeState,
	opts imageparser.ProcessOptions,
) (err error) {
	if opts.Gamma < 0.1 {
		if opts.ApplyColor {
			opts.Gamma = 1.0
		}
	}

	if state.Img == nil {
		return imageparser.ErrNoImageProvided
	}

	// Ensure the source image is in the correct format
	src, _ := gocv.ImageToMatRGBA(state.Img)
	defer src.Close()

	// Perform auto-contrast using Stylization
	result, err := step.AutoContrast(src)
	if err != nil {
		return err
	}
	defer result.Close()

	// Update the pipeline state with the processed image
	state.Img, _ = result.ToImage()
	return nil
}
