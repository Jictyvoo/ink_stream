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

func (step StepRescaleImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	bounds := image.Rect(0, 0, int(step.resolution.Width), int(step.resolution.Height))
	resized := createDrawImage(state.img, bounds)

	drawInterpolator := draw.ApproxBiLinear
	if step.isPixelArt {
		drawInterpolator = draw.NearestNeighbor
	}
	drawInterpolator.Scale(resized, resized.Bounds(), state.img, state.img.Bounds(), draw.Over, nil)

	state.img = resized
	return
}
