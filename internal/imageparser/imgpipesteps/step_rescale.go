package imgpipesteps

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
)

var _ imageparser.PipeStep = (*StepRescaleImage)(nil)

type StepRescaleImage struct {
	resolution deviceprof.Resolution
	isPixelArt bool
	imageparser.BaseImageStep
}

func NewStepRescale(resolution deviceprof.Resolution) *StepRescaleImage {
	return &StepRescaleImage{resolution: resolution}
}

func NewStepThumbnail() StepRescaleImage {
	return StepRescaleImage{resolution: deviceprof.Resolution{Width: 300, Height: 470}}
}

func (step StepRescaleImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	bounds := image.Rect(0, 0, int(step.resolution.Width), int(step.resolution.Height))
	resized := step.DrawImage(state.Img, bounds)

	drawInterpolator := draw.ApproxBiLinear
	if step.isPixelArt {
		drawInterpolator = draw.NearestNeighbor
	}

	// TODO: Check aspect ration to include expand dimensions if necessary
	drawInterpolator.Scale(resized, resized.Bounds(), state.Img, state.Img.Bounds(), draw.Over, nil)

	state.Img = resized
	return
}

func (step StepRescaleImage) calculateNewDimensions(
	bounds image.Rectangle,
) (margins struct{ w, h uint }) {
	actualWidth := uint(bounds.Dx())
	actualHeight := uint(bounds.Dy())

	desiredWidth := step.resolution.Width
	desiredHeight := step.resolution.Height

	// Take some edge cases before floating calculation
	{
		if actualWidth == 0 && actualHeight == 0 {
			margins.w = desiredWidth
			margins.h = desiredHeight
			return
		}
	}

	actualAspect := float64(actualWidth) / float64(actualHeight)
	desiredAspect := float64(desiredWidth) / float64(desiredHeight)

	if actualAspect > desiredAspect {
		// Image is too wide, need to add height (top/bottom padding)
		newHeight := float64(actualWidth) / desiredAspect
		margins.h = uint(newHeight - float64(actualHeight))
	} else if actualAspect < desiredAspect {
		// Image is too tall, need to add width (left/right padding)
		newWidth := float64(actualHeight) * desiredAspect
		margins.w = uint(newWidth - float64(actualWidth))
	}
	// else: aspect ratios match; no margin needed

	return
}
