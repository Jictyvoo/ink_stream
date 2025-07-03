package imgutils

import (
	"image"
	"image/draw"
)

type ImageOrientation uint8

const (
	OrientationPortrait ImageOrientation = iota
	OrientationLandscape
)

func NewOrientation(imgBounds image.Rectangle) ImageOrientation {
	imgOrientation := OrientationPortrait
	if imgBounds.Dx() > imgBounds.Dy() {
		imgOrientation = OrientationLandscape
	}

	return imgOrientation
}

type RotationDegrees uint8

const (
	Rotation90Degrees RotationDegrees = iota
	Rotation180Degrees
	Rotation270Degrees
)

// RotateImage rotates the given image 90 degrees clockwise.
func RotateImage(img image.Image, degrees RotationDegrees) image.Image {
	bounds := img.Bounds()
	var rotated draw.Image

	switch degrees {
	case Rotation90Degrees, Rotation270Degrees:
		// Swap width and height for 90 or 270 degrees
		rotated = NewDrawFromImgColorModel(
			img.ColorModel(),
			image.Rect(0, 0, bounds.Dy(), bounds.Dx()),
		)
	default:
		// Keep the same dimensions for 180 degrees
		rotated = NewDrawFromImgColorModel(img.ColorModel(), bounds)
	}

	// Rotate each pixel
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			switch degrees {
			case Rotation90Degrees:
				rotated.Set(bounds.Max.Y-y-1, x, img.At(x, y))
			case Rotation180Degrees:
				rotated.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
			case Rotation270Degrees:
				rotated.Set(y, bounds.Max.X-x-1, img.At(x, y))
			default:
				// No rotation, return a copy of the original
				rotated.Set(x, y, img.At(x, y))
			}
		}
	}

	return rotated
}
