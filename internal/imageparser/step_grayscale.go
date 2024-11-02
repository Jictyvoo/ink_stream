package imageparser

import (
	"image"
	"image/color"
)

type StepGrayScaleImage struct{}

func (sgsi StepGrayScaleImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	grayImg := image.NewGray(state.img.Bounds())
	for y := 0; y < state.img.Bounds().Dy(); y++ {
		for x := 0; x < state.img.Bounds().Dx(); x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(state.img.At(x, y)))
		}
	}

	state.img = grayImg
	return
}
