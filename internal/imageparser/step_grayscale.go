package imageparser

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var (
	_ UnitStep = (*StepGrayScaleImage)(nil)
	_ PipeStep = (*StepGrayScaleImage)(nil)
)

type StepGrayScaleImage struct {
	baseImageStep
}

func NewStepGrayScale() *StepGrayScaleImage {
	return &StepGrayScaleImage{}
}

func (step StepGrayScaleImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	grayImg := image.NewGray(state.img.Bounds())
	for x, y := range imgutils.Iterator(state.img) {
		grayImg.Set(x, y, step.PixelStep(state.img.At(x, y)))
	}

	state.img = grayImg
	return
}

func (step StepGrayScaleImage) PixelStep(imgColor color.Color) color.Color {
	return color.GrayModel.Convert(imgColor)
}
