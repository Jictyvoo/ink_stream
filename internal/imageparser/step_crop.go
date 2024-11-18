package imageparser

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ PipeStep = (*StepCropRotateImage)(nil)

type StepCropRotateImage struct {
	palette     imgutils.ColorConverter
	rotateImage bool
	orientation imgutils.ImageOrientation
	baseImageStep
}

func NewStepCropRotate(
	rotate bool, palette color.Palette,
	orientation imgutils.ImageOrientation,
) *StepCropRotateImage {
	return &StepCropRotateImage{
		palette:       palette,
		rotateImage:   rotate,
		orientation:   orientation,
		baseImageStep: baseImageStep{fac: imgutils.NewImageFactory(palette)},
	}
}

func (step StepCropRotateImage) PerformExec(state *pipeState, _ processOptions) (err error) {
	// Step 1: Crop the image to exclude unnecessary parts and add margins
	desiredBox := state.img.Bounds()
	// TODO: Include gaussian blur
	croppedBox := imgutils.CropBox(state.img, step.palette, imgutils.BoxEliminateMinimumColor)
	if croppedBox != desiredBox {
		desiredBox = croppedBox
	}

	// Step 2: Check dimensions and apply logic based on width and height
	state.img = step.cropImage(state.img, desiredBox)
	newBounds := state.img.Bounds()
	if step.orientation == imgutils.OrientationPortrait && newBounds.Dx() > newBounds.Dy() {
		if step.rotateImage {
			// Rotate the image if rotateImage is true
			state.img = imgutils.RotateImage(state.img, imgutils.Rotation90Degrees)
		} else {
			// Cut the image in half horizontally
			midX := newBounds.Min.X + newBounds.Dx()/2
			halfBounds := struct{ left, right image.Rectangle }{
				left:  image.Rect(newBounds.Min.X, newBounds.Min.Y, midX, newBounds.Max.Y),
				right: image.Rect(midX, newBounds.Min.Y, newBounds.Max.X, newBounds.Max.Y),
			}

			originalImg := state.img
			state.img = step.cropImage(originalImg, halfBounds.left)
			state.img = step.cropImage(originalImg, halfBounds.right)
		}
	}

	return
}

// Helper function: Crop an image to the given rectangle
func (step StepCropRotateImage) cropImage(img image.Image, rect image.Rectangle) image.Image {
	cropped := imgutils.NewDrawFromImgColorModel(img, rect)
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)
	return cropped
}
