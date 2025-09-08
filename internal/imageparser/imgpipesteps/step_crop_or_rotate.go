package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

var _ imageparser.PipeStep = (*StepCropOrRotateImage)(nil)

type StepCropOrRotateImage struct {
	rotateImage bool
	orientation inktypes.ImageOrientation
	imageparser.BaseImageStep
}

func NewStepCropOrRotate(
	rotate bool, palette color.Palette,
	orientation inktypes.ImageOrientation,
) *StepCropOrRotateImage {
	return &StepCropOrRotateImage{
		rotateImage:   rotate,
		orientation:   orientation,
		BaseImageStep: imageparser.NewBaseImageStep(palette),
	}
}

func (step StepCropOrRotateImage) StepID() string {
	return "crop_or_rotate"
}

func (step StepCropOrRotateImage) PerformExec(
	state *imageparser.PipeState, _ imageparser.ProcessOptions,
) (err error) {
	originalBounds := state.Img.Bounds()
	imgOrientation := imgutils.NewOrientation(originalBounds)

	if step.orientation != imgOrientation {
		if step.rotateImage {
			// Rotate the image if rotateImage is true
			state.Img = imgutils.RotateImage(state.Img, imgutils.Rotation90Degrees)
		} else {
			// Cut the image in half horizontally
			halfBounds := imgutils.HalfSplit(originalBounds, imgOrientation)

			originalImg := state.Img
			switch step.orientation {
			case inktypes.OrientationLandscape:
				state.SubImages = []image.Image{
					imgutils.CropImage(originalImg, halfBounds.Top),
					imgutils.CropImage(originalImg, halfBounds.Bottom),
				}
			case inktypes.OrientationPortrait:
				state.SubImages = []image.Image{
					imgutils.CropImage(originalImg, halfBounds.Left),
					imgutils.CropImage(originalImg, halfBounds.Right),
				}
			}
		}
	}

	return err
}
