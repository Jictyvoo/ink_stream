//go:build cgo

package imgpipesteps

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
)

var _ imageparser.PipeStep = (*StepApplyGaussianBlurImage)(nil)

type (
	gaussianKernel struct {
		radius int
		sigma  float64
		kSize  image.Point
	}
	StepApplyGaussianBlurImage struct {
		kernel gaussianKernel
		imageparser.BaseImageStep
	}
)

func NewStepGaussianBlur(radius int) *StepApplyGaussianBlurImage {
	return &StepApplyGaussianBlurImage{kernel: createBlurKernel(radius)}
}

func createBlurKernel(radius int) gaussianKernel {
	ksize := image.Point{X: (2 * radius) + 1, Y: (2 * radius) + 1}

	// Sigma values: heuristic based on radius
	sigma := float64(radius) / 2.0
	return gaussianKernel{radius: radius, sigma: sigma, kSize: ksize}
}

func (step StepApplyGaussianBlurImage) PerformExec(
	state *imageparser.PipeState, _ imageparser.ProcessOptions,
) error {
	if step.kernel.radius <= 0 || step.kernel.radius > 200 {
		return nil
	}

	img := state.Img
	rgbaMat, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		return err
	}
	defer rgbaMat.Close()

	// Create a destination Mat to store the blurred image
	dst := gocv.NewMat()
	defer dst.Close()

	// Apply GaussianBlur
	gocv.GaussianBlur(
		rgbaMat, &dst,
		step.kernel.kSize,
		step.kernel.sigma,
		step.kernel.sigma,
		gocv.BorderDefault,
	)
	state.Img, err = dst.ToImage()
	return err
}
