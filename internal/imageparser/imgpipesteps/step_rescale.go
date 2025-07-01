package imgpipesteps

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
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
	paddedImage := step.wrapImgWithMargins(state.Img)

	// Now resize the padded image to target resolution
	bounds := image.Rect(0, 0, int(step.resolution.Width), int(step.resolution.Height))
	resized := step.DrawImage(state.Img.ColorModel(), bounds)

	drawInterpolator := draw.ApproxBiLinear
	if step.isPixelArt {
		drawInterpolator = draw.NearestNeighbor
	}

	drawInterpolator.Scale(
		resized, resized.Bounds(),
		paddedImage, paddedImage.Bounds(),
		draw.Over, nil,
	)

	state.Img = resized
	return
}

func (step StepRescaleImage) wrapImgWithMargins(sttImg image.Image) image.Image {
	// Calculate required padding
	margins := step.calculateNewDimensions(sttImg.Bounds())
	if margins.w == 0 && margins.h == 0 {
		return sttImg
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
	{
		dominantColors := imgutils.ImageMarginDominantColor(
			sttImg, uint32(margins.w), uint32(margins.h), 5,
		)
		// Fill top and bottom padding
		if offsets.Y > 0 {
			topRect := image.Rect(0, 0, paddedWidth, offsets.Y)
			imgutils.FillImageRegionWithColor(paddedImage, topRect, dominantColors.Top)

			bottomRect := image.Rect(0, imgBounds.Dy()+offsets.Y, paddedWidth, paddedHeight)
			imgutils.FillImageRegionWithColor(paddedImage, bottomRect, dominantColors.Bottom)
		}

		// Fill left and right padding
		if offsets.X > 0 {
			leftRect := image.Rect(0, 0, offsets.X, paddedHeight)
			imgutils.FillImageRegionWithColor(paddedImage, leftRect, dominantColors.Left)

			rightRect := image.Rect(imgBounds.Dx()+offsets.X, 0, paddedWidth, paddedHeight)
			imgutils.FillImageRegionWithColor(paddedImage, rightRect, dominantColors.Right)
		}
	}

	// Draw the original image onto the center of the padded image
	draw.Draw(
		paddedImage, image.Rect(
			offsets.X, offsets.Y,
			offsets.X+imgBounds.Dx(), offsets.Y+imgBounds.Dy(),
		),
		sttImg, image.Point{}, draw.Src,
	)
	return paddedImage
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
