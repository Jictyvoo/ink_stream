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
	blurSubstep *StepApplyGaussianBlurImage
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
		blurSubstep:   NewStepGaussianBlur(5),
		BaseImageStep: imageparser.NewBaseImageStep(palette),
	}
}

func (step StepCropRotateImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	// Step 1: Crop the image to exclude unnecessary parts and add margins
	desiredBox := state.Img.Bounds()
	originalImage := state.Img
	// Use gaussian blur to perform image crop
	if err = step.blurSubstep.PerformExec(state, imageparser.ProcessOptions{}); err != nil {
		state.Img = originalImage
		return
	}
	croppedBox := imgutils.CropBox(state.Img, step.palette, imgutils.BoxEliminateMinimumColor)
	if croppedBox != desiredBox {
		// Prevent cropping if new dimensions are lower than 80% the original size
		originalSize := [2]int{desiredBox.Dx(), desiredBox.Dy()}
		newSize := [2]int{croppedBox.Dx(), croppedBox.Dy()}
		if (newSize[0]*100)/originalSize[0] >= 80 && (newSize[1]*100)/originalSize[1] >= 80 {
			desiredBox = croppedBox
		}
	}

	// Step 2: Check dimensions and apply logic based on width and height
	state.Img = step.cropImage(originalImage, desiredBox)
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
