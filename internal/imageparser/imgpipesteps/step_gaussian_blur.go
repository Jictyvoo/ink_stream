package imgpipesteps

import (
	"image"
	"image/color"
	"math"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepApplyGaussianBlurImage)(nil)

type (
	gaussianKernel struct {
		radius    int
		matrix    [][]float64
		weightSum float64
	}
	StepApplyGaussianBlurImage struct {
		kernel gaussianKernel
		imageparser.BaseImageStep
	}
)

func NewStepGaussianBlur(radius int) *StepApplyGaussianBlurImage {
	var gk gaussianKernel
	gk.initKernel(radius)

	return &StepApplyGaussianBlurImage{kernel: gk}
}

func (gk *gaussianKernel) initKernel(radius int) (kernel [][]float64, kSum float64) {
	gk.radius = radius
	kernelWidth := (2 * gk.radius) + 1
	kernel = make([][]float64, kernelWidth)
	for index := range kernelWidth {
		kernel[index] = make([]float64, kernelWidth)
	}

	// gaussKernelPoint calculates the weight of a pixel based on its distance
	gaussKernelPoint := func(distanceSquared float64) float64 {
		if distanceSquared < 0 {
			return 0
		}
		sigma := max(float64(gk.radius)/2, 1)
		exponentDenominator := 2 * sigma * sigma
		return math.Exp(-distanceSquared/(exponentDenominator)) / (2 * math.Pi * sigma * sigma)
	}

	for dy := -gk.radius; dy <= gk.radius; dy++ {
		for dx := -gk.radius; dx <= gk.radius; dx++ {
			distanceSquared := float64((dx * dx) + (dy * dy))
			kernelValue := gaussKernelPoint(distanceSquared)
			kernel[dx+gk.radius][dy+gk.radius] = kernelValue
			kSum += kernelValue
		}
	}
	gk.matrix = kernel
	gk.weightSum = kSum

	return
}

func (gk *gaussianKernel) onPoint(x, y int) float64 {
	return gk.matrix[x][y] / gk.weightSum
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
	for x, y := range imgutils.Iterator(img) {
		blurredImg.Set(x, y, step.NeighborCalculation(img, x, y))
	}

	// Update the state with the blurred image
	state.Img = blurredImg
	return
}

// NeighborCalculation calculates the blurred value of a pixel using Gaussian weights
func (step StepApplyGaussianBlurImage) NeighborCalculation(
	img image.Image, x, y int,
) color.Color {
	var colorSum struct{ red, green, blue, alpha float64 }
	radius := step.kernel.radius
	imgBounds := img.Bounds()

	// This is the convolution step
	// Run the kernel over this grouping of pixels centered around the pixel at (x,y)
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			px, py := x-dx, y-dy
			// Ignore loop if is out of bounds
			if px < imgBounds.Min.X || py < imgBounds.Min.Y ||
				px > imgBounds.Max.X || py > imgBounds.Max.Y {
				continue
			}

			// Get the color values of the current pixel
			r, g, b, a := img.At(px, py).RGBA()
			kernelValue := step.kernel.onPoint(dx+radius, dy+radius)

			// Accumulate weighted color values
			colorSum.red += float64(r>>8) * kernelValue
			colorSum.green += float64(g>>8) * kernelValue
			colorSum.blue += float64(b>>8) * kernelValue
			colorSum.alpha += float64(a>>8) * kernelValue
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
