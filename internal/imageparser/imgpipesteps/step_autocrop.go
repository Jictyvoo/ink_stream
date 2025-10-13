package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepAutoCropImage)(nil)

type StepAutoCropImage struct {
	palette     imgutils.ColorConverter
	blurSubstep *StepApplyGaussianBlurImage
	imageparser.BaseImageStep
}

func NewStepAutoCrop(
	palette color.Palette,
) *StepAutoCropImage {
	return &StepAutoCropImage{
		palette:       palette,
		blurSubstep:   NewStepGaussianBlur(5),
		BaseImageStep: imageparser.NewBaseImageStep(palette),
	}
}

func (step StepAutoCropImage) StepID() string {
	return "autocrop"
}

func (step StepAutoCropImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	// Step 1: Crop the image to exclude unnecessary parts and add margins
	originalImage := state.Img
	originalBox := originalImage.Bounds()
	desiredBox := originalBox
	// Use gaussian blur to perform image crop
	if err = step.blurSubstep.PerformExec(state, imageparser.ProcessOptions{}); err != nil {
		state.Img = originalImage
		return err
	}

	blurredImg := state.Img
	state.Img = originalImage
	croppedBox := imgutils.CropBox(blurredImg, step.palette, imgutils.BoxEliminateMinimumColor)
	if croppedBox != desiredBox {
		// Prevent cropping if new dimensions are lower than 80% the original size
		originalSize := [2]int{desiredBox.Dx(), desiredBox.Dy()}
		newSize := [2]int{croppedBox.Dx(), croppedBox.Dy()}
		if (newSize[0]*100)/originalSize[0] >= 80 && (newSize[1]*100)/originalSize[1] >= 80 {
			// Include a little margin on croppedBox
			desiredBox = step.wrapInMargin(croppedBox, originalBox.Max)
		}
	}

	// Step 2: Check dimensions and apply logic based on width and height
	if desiredBox != originalBox {
		state.Img = imgutils.CropImage(state.Img, desiredBox)
	}

	return err
}

func (step StepAutoCropImage) wrapInMargin(
	croppedBox image.Rectangle, limits image.Point,
) image.Rectangle {
	desiredBox := imgutils.MarginBox(croppedBox, 1.6e-2)
	if desiredBox.Max.X > limits.X {
		desiredBox.Max.X = limits.X
	}
	if desiredBox.Max.Y > limits.Y {
		desiredBox.Max.Y = limits.Y
	}
	return desiredBox
}
