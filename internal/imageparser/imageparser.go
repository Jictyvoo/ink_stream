package imageparser

import (
	"image"
	"image/color"
	"image/draw"
)

type (
	UnitStep interface {
		PixelStep(imgColor color.Color) color.Color
	}

	PipeStep interface {
		PerformExec(state *pipeState, opts processOptions) (err error)
	}
)

func createDrawImage(img image.Image, bounds image.Rectangle) draw.Image {
	switch img.ColorModel() {
	case color.GrayModel, color.Gray16Model:
		return image.NewGray(bounds)
	case color.RGBAModel, color.RGBA64Model, color.NRGBAModel, color.NRGBA64Model:
		return image.NewRGBA(bounds)
	}

	return image.NewRGBA(bounds)
}
