package imgpipesteps

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/comic_manga-extractor-Converter/internal/imageparser"
	"github.com/Jictyvoo/comic_manga-extractor-Converter/pkg/deviceprof"
)

var _ imageparser.PipeStep = (*StepRescaleImage)(nil)

type StepRescaleImage struct {
	resolution deviceprof.Resolution
	isPixelArt bool
	imageparser.BaseImageStep
}

func NewStepRescale(resolution deviceprof.Resolution, allowStretch bool) *StepRescaleImage {
	return &StepRescaleImage{
		resolution: resolution,
	}
}

func NewStepThumbnail() StepRescaleImage {
	return StepRescaleImage{resolution: deviceprof.Resolution{Width: 300, Height: 470}}
}

func (step StepRescaleImage) StepID() string {
	return "rescale"
}

func (step StepRescaleImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	inputImage := state.Img
	// Now resize the padded image to target resolution
	bounds := image.Rect(0, 0, int(step.resolution.Width), int(step.resolution.Height))
	resized := step.DrawImage(state.Img.ColorModel(), bounds)

	drawInterpolator := draw.ApproxBiLinear
	if step.isPixelArt {
		drawInterpolator = draw.NearestNeighbor
	}

	drawInterpolator.Scale(
		resized, resized.Bounds(),
		inputImage, inputImage.Bounds(),
		draw.Over, nil,
	)

	state.Img = resized
	return
}
