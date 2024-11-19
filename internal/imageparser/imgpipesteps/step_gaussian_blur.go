package imgpipesteps

import (
	"image"
	"image/color"
	"math"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepApplyGaussianBlurImage)(nil)

type StepApplyGaussianBlurImage struct {
	palette imgutils.ColorConverter
	radius  int
	imageparser.BaseImageStep
}

func NewStepGaussianBlur(radius int) *StepApplyGaussianBlurImage {
	return &StepApplyGaussianBlurImage{radius: radius}
}

// gaussKernelPoint calculates the weight of a pixel based on its distance
func (step StepApplyGaussianBlurImage) gaussKernelPoint(distanceSquared float64) float64 {
	if distanceSquared < 0 {
		return 0
	}
	sigma := max(float64(step.radius)/2, 1)
	exponentDenominator := 2 * sigma * sigma
	return math.Exp(-distanceSquared/(exponentDenominator)) / (2 * math.Pi * sigma * sigma)
}

func (step StepApplyGaussianBlurImage) makeKernel() (kernel [][]float64, kSum float64) {
	kernelWidth := (2 * step.radius) + 1
	kernel = make([][]float64, kernelWidth)
	for index := range kernelWidth {
		kernel[index] = make([]float64, kernelWidth)
	}

	for dy := -step.radius; dy <= step.radius; dy++ {
		for dx := -step.radius; dx <= step.radius; dx++ {
			distanceSquared := float64((dx * dx) + (dy * dy))
			kernelValue := step.gaussKernelPoint(distanceSquared)
			kernel[dx+step.radius][dy+step.radius] = kernelValue
			kSum += kernelValue
		}
	}

	return
}

// PerformExec applies the Gaussian blur to the entire image
func (step StepApplyGaussianBlurImage) PerformExec(
	state *imageparser.PipeState, _ imageparser.ProcessOptions,
) (err error) {
	img := state.Img

	// Create a new image with the same size as the input image
	bounds := img.Bounds()
	blurredImg := step.DrawImage(img, bounds)

	// Iterate through each pixel and apply the Gaussian blur
	kernel, kSum := step.makeKernel()
	for x, y := range imgutils.Iterator(img) {
		blurredImg.Set(x, y, step.NeighborCalculation(img, x, y, kernel, kSum))
	}

	// Update the state with the blurred image
	state.Img = blurredImg
	return
}

// NeighborCalculation calculates the blurred value of a pixel using Gaussian weights
func (step StepApplyGaussianBlurImage) NeighborCalculation(
	img image.Image, x, y int,
	kernel [][]float64, kSum float64,
) color.Color {
	var colorSum struct{ red, green, blue, alpha float64 }

	// Bounds of the image
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Apply the Gaussian kernel to all pixels in the neighborhood
	for dy := -step.radius; dy <= step.radius; dy++ {
		for dx := -step.radius; dx <= step.radius; dx++ {
			px, py := x-dx, y-dy
			if px >= 0 && px < w && py >= 0 && py < h {
				// Get the color values of the current pixel
				r, g, b, a := img.At(px, py).RGBA()
				kernelValue := kernel[dx+step.radius][dy+step.radius] / kSum

				// Accumulate weighted color values
				colorSum.red += float64(r>>8) * kernelValue
				colorSum.green += float64(g>>8) * kernelValue
				colorSum.blue += float64(b>>8) * kernelValue
				colorSum.alpha += float64(a>>8) * kernelValue
			}
		}
	}

	// Normalize the accumulated values
	return color.NRGBA{
		R: imgutils.NormalizePixel(colorSum.red),
		G: imgutils.NormalizePixel(colorSum.green),
		B: imgutils.NormalizePixel(colorSum.blue),
		A: imgutils.NormalizePixel(colorSum.alpha),
	}
}
