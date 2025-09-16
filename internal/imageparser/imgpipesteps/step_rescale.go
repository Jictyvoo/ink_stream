package imgpipesteps

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

var _ imageparser.PipeStep = (*StepRescaleImage)(nil)

type StepRescaleImage struct {
	resolution   inktypes.ImageDimensions
	isPixelArt   bool
	allowStretch bool
	imageparser.BaseImageStep
}

func NewStepRescale(resolution inktypes.ImageDimensions, allowStretch bool) *StepRescaleImage {
	return &StepRescaleImage{
		resolution:   resolution,
		allowStretch: allowStretch,
	}
}

func NewStepThumbnail() StepRescaleImage {
	return StepRescaleImage{
		allowStretch: true,
		resolution:   inktypes.ImageDimensions{Width: 300, Height: 470},
	}
}

func (step StepRescaleImage) StepID() string {
	return "rescale"
}

func (step StepRescaleImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	inputImage := state.Img
	if !step.allowStretch {
		step.resolution = step.updateTargetResolution(
			inktypes.ImageDimensions{
				Width:  uint16(inputImage.Bounds().Dx()),
				Height: uint16(inputImage.Bounds().Dy()),
			},
		)
	}

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
	return err
}

func (step StepRescaleImage) updateTargetResolution(
	imgDimensions inktypes.ImageDimensions,
) inktypes.ImageDimensions {
	actualAspect := float64(imgDimensions.Width) / float64(imgDimensions.Height)
	desiredAspect := float64(step.resolution.Width) / float64(step.resolution.Height)

	newWidth, newHeight := float64(step.resolution.Width), float64(step.resolution.Height)
	if actualAspect > desiredAspect {
		// Image has a bigger width than target, need to adjust height
		widthProportion := float64(step.resolution.Width) / float64(imgDimensions.Width)
		newHeight = float64(imgDimensions.Height) * widthProportion
	} else if actualAspect < desiredAspect {
		// Image has a bigger height than target, need to adjust width
		heightProportion := float64(step.resolution.Height) / float64(imgDimensions.Height)
		newWidth = float64(imgDimensions.Width) * heightProportion
	}

	return inktypes.ImageDimensions{Width: uint16(newWidth), Height: uint16(newHeight)}
}
