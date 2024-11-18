package imgpipesteps

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepCropRotateImage)(nil)

type StepCropRotateImage struct {
	palette     imgutils.ColorConverter
	rotateImage bool
	orientation imgutils.ImageOrientation
	imageparser.BaseImageStep
}

func NewStepCropRotate(
	rotate bool, palette color.Palette,
	orientation imgutils.ImageOrientation,
) *StepCropRotateImage {
	return &StepCropRotateImage{
		palette:       palette,
		rotateImage:   rotate,
		orientation:   orientation,
		BaseImageStep: imageparser.NewBaseImageStep(palette),
	}
}

func (step StepCropRotateImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	// Step 1: Crop the image to exclude unnecessary parts and add margins
	desiredBox := state.Img.Bounds()
	// TODO: Include gaussian blur
	croppedBox := imgutils.CropBox(state.Img, step.palette, imgutils.BoxEliminateMinimumColor)
	if croppedBox != desiredBox {
		desiredBox = croppedBox
	}

	// Step 2: Check dimensions and apply logic based on width and height
	state.Img = step.cropImage(state.Img, desiredBox)
	newBounds := state.Img.Bounds()
	if step.orientation == imgutils.OrientationPortrait && newBounds.Dx() > newBounds.Dy() {
		if step.rotateImage {
			// Rotate the image if rotateImage is true
			state.Img = imgutils.RotateImage(state.Img, imgutils.Rotation90Degrees)
		} else {
			// Cut the image in half horizontally
			midX := newBounds.Min.X + newBounds.Dx()/2
			halfBounds := struct{ left, right image.Rectangle }{
				left:  image.Rect(newBounds.Min.X, newBounds.Min.Y, midX, newBounds.Max.Y),
				right: image.Rect(midX, newBounds.Min.Y, newBounds.Max.X, newBounds.Max.Y),
			}

			originalImg := state.Img
			state.Img = step.cropImage(originalImg, halfBounds.left)
			state.Img = step.cropImage(originalImg, halfBounds.right)
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
