package imageparser

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/ink_stream/internal/deviceprof"
)

type StepRescaleImage struct {
	resolution deviceprof.Resolution
	isPixelArt bool
}

func NewStepRescale(resolution deviceprof.Resolution) StepRescaleImage {
	return StepRescaleImage{resolution: resolution}
}

func NewStepThumbnail() StepRescaleImage {
	return StepRescaleImage{resolution: deviceprof.Resolution{Width: 300, Height: 470}}
}

func (sgsi StepRescaleImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	resized := image.NewRGBA(image.Rect(0, 0, int(sgsi.resolution.Width), int(sgsi.resolution.Height)))

	drawInterpolator := draw.ApproxBiLinear
	if sgsi.isPixelArt {
		drawInterpolator = draw.NearestNeighbor
	}
	drawInterpolator.Scale(resized, resized.Bounds(), state.img, state.img.Bounds(), draw.Over, nil)

	state.img = resized
	return
}
