//go:build !cgo

package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepApplyGaussianBlurImage)(nil)

type (
	gaussianKernel struct {
		radius  int
		weights []int
	}
	StepApplyGaussianBlurImage struct {
		kernel gaussianKernel
		imageparser.BaseImageStep
	}
)

func NewStepGaussianBlur(radius int) *StepApplyGaussianBlurImage {
	return &StepApplyGaussianBlurImage{kernel: createBlurKernel(radius)}
}

func (step StepApplyGaussianBlurImage) StepID() string {
	return "gaussian_blur"
}

func createBlurKernel(radius int) gaussianKernel {
	size := (2 * radius) + 1
	weights := make([]int, size)
	for i := 0; i <= radius; i++ {
		weights[i] = 16 * (i + 1)
		weights[size-1-i] = weights[i]
	}
	return gaussianKernel{radius: radius, weights: weights}
}

func (step StepApplyGaussianBlurImage) PerformExec(
	state *imageparser.PipeState, _ imageparser.ProcessOptions,
) error {
	if step.kernel.radius <= 0 || step.kernel.radius > 200 {
		return nil
	}

	img := state.Img
	bounds := img.Bounds()
	blurredImg := step.DrawImage(img.ColorModel(), bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newColor := step.applyKernel(img, x, y, bounds)
			blurredImg.Set(x, y, newColor)
		}
	}

	state.Img = blurredImg
	return nil
}

func (step StepApplyGaussianBlurImage) applyKernel(
	img image.Image, x, y int, bounds image.Rectangle,
) color.Color {
	var colorSums struct{ r, g, b, a, wc, wa int64 }

	for wy := 0; wy < len(step.kernel.weights); wy++ {
		srcY := y + wy - step.kernel.radius
		if srcY < bounds.Min.Y || srcY >= bounds.Max.Y {
			continue
		}
		for wx := 0; wx < len(step.kernel.weights); wx++ {
			srcX := x + wx - step.kernel.radius
			if srcX < bounds.Min.X || srcX >= bounds.Max.X {
				continue
			}

			c := img.At(srcX, srcY)
			r, g, b, a := c.RGBA()
			wp := int64(step.kernel.weights[wx] * step.kernel.weights[wy])

			colorSums.wa += wp
			wp *= int64(a) + int64(a>>7)
			colorSums.wc += wp
			wp >>= 8

			colorSums.a += wp * int64(a)
			colorSums.r += wp * int64(r)
			colorSums.g += wp * int64(g)
			colorSums.b += wp * int64(b)
		}
	}

	if colorSums.wa == 0 || colorSums.wc == 0 {
		return img.At(x, y) // Fallback to original color
	}

	return color.RGBA{
		A: imgutils.NormalizePixel(colorSums.a / colorSums.wa),
		R: imgutils.NormalizePixel(colorSums.r / colorSums.wc),
		G: imgutils.NormalizePixel(colorSums.g / colorSums.wc),
		B: imgutils.NormalizePixel(colorSums.b / colorSums.wc),
	}
}
