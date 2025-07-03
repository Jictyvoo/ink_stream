package imgpipesteps

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

var _ imageparser.PipeStep = (*StepMarginWrapImage)(nil)

type StepMarginWrapImage struct {
	resolution  deviceprof.Resolution
	marginColor color.Color
	imageparser.BaseImageStep
}

func NewStepMarginWrap(resolution deviceprof.Resolution) *StepMarginWrapImage {
	return &StepMarginWrapImage{
		resolution:  resolution,
		marginColor: color.White,
	}
}

func (step StepMarginWrapImage) StepID() string {
	return "wrap_in_margin"
}

func (step StepMarginWrapImage) PerformExec(
	state *imageparser.PipeState,
	_ imageparser.ProcessOptions,
) (err error) {
	// Calculate required padding
	sttImg := state.Img
	margins := step.calculateNewDimensions(sttImg.Bounds())
	if margins.w == 0 && margins.h == 0 {
		return
	}

	imgBounds := sttImg.Bounds()
	// Apply padding to the image (centered)
	paddedWidth := imgBounds.Dx() + int(margins.w)
	paddedHeight := imgBounds.Dy() + int(margins.h)

	offsets := image.Point{
		X: int(margins.w) / 2,
		Y: int(margins.h) / 2,
	}
	paddedImage := image.NewNRGBA(image.Rect(0, 0, paddedWidth, paddedHeight))

	// Fill image paddings
	marginColors := imgutils.Margins[color.Color]{
		Top:    step.marginColor,
		Bottom: step.marginColor,
		Left:   step.marginColor,
		Right:  step.marginColor,
	}
	marginColors.UpdateNonEmpty(imgutils.ImageMarginDominantColor(
		sttImg, uint32(margins.w), uint32(margins.h), 5,
	))
	drawMargins(paddedImage, offsets, imgBounds, marginColors)

	// Draw the original image onto the center of the padded image
	draw.Draw(
		paddedImage, image.Rect(
			offsets.X, offsets.Y,
			offsets.X+imgBounds.Dx(), offsets.Y+imgBounds.Dy(),
		),
		sttImg, sttImg.Bounds().Min, draw.Src,
	)

	state.Img = paddedImage
	return
}

func drawMargins(
	dstImg draw.Image, offsets image.Point,
	originalBounds image.Rectangle,
	marginColors imgutils.Margins[color.Color],
) {
	imgBounds := dstImg.Bounds()

	// Fill top and bottom padding
	if offsets.Y > 0 {
		topRect := image.Rect(0, 0, imgBounds.Dx(), offsets.Y)
		imgutils.FillImageRegionWithColor(dstImg, topRect, marginColors.Top)

		bottomRect := image.Rect(0, originalBounds.Dy()+offsets.Y, imgBounds.Dx(), imgBounds.Dy())
		imgutils.FillImageRegionWithColor(dstImg, bottomRect, marginColors.Bottom)
	}

	// Fill left and right padding
	if offsets.X > 0 {
		leftRect := image.Rect(0, 0, offsets.X, imgBounds.Dy())
		imgutils.FillImageRegionWithColor(dstImg, leftRect, marginColors.Left)

		rightRect := image.Rect(originalBounds.Dx()+offsets.X, 0, imgBounds.Dx(), imgBounds.Dy())
		imgutils.FillImageRegionWithColor(dstImg, rightRect, marginColors.Right)
	}
}

func (step StepMarginWrapImage) calculateNewDimensions(
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
