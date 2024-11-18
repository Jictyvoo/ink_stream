package imageparser

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

type (
	paletteFactoryStep interface {
		UpdateDrawFactory(fac imgutils.DrawImageFactory)
		privateInternalStep()
	}
	UnitStep interface {
		PixelStep(imgColor color.Color) color.Color
	}

	PipeStep interface {
		PerformExec(state *PipeState, opts ProcessOptions) (err error)
		paletteFactoryStep
	}
)

type BaseImageStep struct {
	fac imgutils.DrawImageFactory
}

func NewBaseImageStep(palette color.Palette) BaseImageStep {
	return BaseImageStep{fac: imgutils.NewImageFactory(palette)}
}

func (s *BaseImageStep) privateInternalStep() {
	// Do nothing, this function only exists to make sure that all types compose this struct
}

func (s *BaseImageStep) DrawImage(img image.Image, bounds image.Rectangle) draw.Image {
	if s.fac != nil {
		return s.fac.CreateDrawImage(img, bounds)
	}
	return imgutils.NewDrawFromImgColorModel(img, bounds)
}

func (s *BaseImageStep) UpdateDrawFactory(fac imgutils.DrawImageFactory) {
	s.fac = fac
}
