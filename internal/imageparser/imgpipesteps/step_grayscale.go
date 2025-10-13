package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var (
	_ imageparser.UnitStep = (*StepGrayScaleImage)(nil)
	_ imageparser.PipeStep = (*StepGrayScaleImage)(nil)
)

type StepGrayScaleImage struct {
	imageparser.BaseImageStep
}

func NewStepGrayScale() *StepGrayScaleImage {
	return &StepGrayScaleImage{}
}

func (step StepGrayScaleImage) StepID() string {
	return "grayscale"
}

func (step StepGrayScaleImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	switch state.Img.ColorModel() { // Prevent redraw in case already is in grayscale
	case color.GrayModel, color.Gray16Model:
		return nil
	}

	grayImg := image.NewGray(state.Img.Bounds())
	for x, y := range imgutils.Iterator(state.Img) {
		grayImg.Set(x, y, step.PixelStep(state.Img.At(x, y)))
	}

	state.Img = grayImg
	return err
}

func (step StepGrayScaleImage) PixelStep(imgColor color.Color) color.Color {
	return color.GrayModel.Convert(imgColor)
}
