package imgpipesteps

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepCropOrRotateImage)(nil)

type StepCropOrRotateImage struct {
	rotateImage bool
	orientation imgutils.ImageOrientation
	imageparser.BaseImageStep
}

func NewStepCropOrRotate(
	rotate bool, palette color.Palette,
	orientation imgutils.ImageOrientation,
) *StepCropOrRotateImage {
	return &StepCropOrRotateImage{
		rotateImage:   rotate,
		orientation:   orientation,
		BaseImageStep: imageparser.NewBaseImageStep(palette),
	}
}

func (step StepCropOrRotateImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	originalBounds := state.Img.Bounds()
	if step.orientation == imgutils.OrientationPortrait &&
		originalBounds.Dx() > originalBounds.Dy() {
		if step.rotateImage {
			// Rotate the image if rotateImage is true
			state.Img = imgutils.RotateImage(state.Img, imgutils.Rotation90Degrees)
		} else {
			// Cut the image in half horizontally
			midX := originalBounds.Min.X + originalBounds.Dx()/2
			halfBounds := struct{ left, right image.Rectangle }{
				left:  image.Rect(originalBounds.Min.X, originalBounds.Min.Y, midX, originalBounds.Max.Y),
				right: image.Rect(midX, originalBounds.Min.Y, originalBounds.Max.X, originalBounds.Max.Y),
			}

			originalImg := state.Img
			state.SubImages = []image.Image{
				imgutils.CropImage(originalImg, halfBounds.left),
				imgutils.CropImage(originalImg, halfBounds.right),
			}
		}
	}

	return
}
