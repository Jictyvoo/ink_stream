package imageparser

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils"
)

type StepGrayScaleImage struct{}

func NewStepGrayScale() StepGrayScaleImage {
	return StepGrayScaleImage{}
}

func (step StepGrayScaleImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	grayImg := image.NewGray(state.img.Bounds())
	for x, y := range imgutils.Iterator(state.img) {
		grayImg.Set(x, y, color.GrayModel.Convert(state.img.At(x, y)))
	}

	state.img = grayImg
	return
}
