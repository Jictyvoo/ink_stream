package imgutils

import (
	"image"
	"image/color"
	"image/draw"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

const (
	MaxPixelValue = (1 << 8) - 1
	MinPixelValue = 0
)

func NormalizePixel[T Number](value T) uint8 {
	// The conversion to int64 is required because the maximum value for an int8 is 127
	return uint8(max(MinPixelValue, min(MaxPixelValue, int64(value))))
}

type (
	ColorConverter interface {
		Convert(c color.Color) color.Color
	}
	DrawImageFactory interface {
		CreateDrawImage(colorModel color.Model, bounds image.Rectangle) draw.Image
	}
)

type Margins[T comparable] struct {
	Top, Bottom, Left, Right T
}

func (c *Margins[T]) UpdateNonEmpty(other Margins[T]) {
	var defaultValue T
	if c.Left == defaultValue {
		c.Left = other.Left
	}
	if c.Right == defaultValue {
		c.Right = other.Right
	}
	if c.Bottom == defaultValue {
		c.Bottom = other.Bottom
	}
	if c.Top == defaultValue {
		c.Top = other.Top
	}
}

func FillImageRegionWithColor(img draw.Image, region image.Rectangle, col color.Color) {
	draw.Draw(img, region, &image.Uniform{C: col}, image.Point{}, draw.Src)
}

func CropImage(img image.Image, rect image.Rectangle) image.Image {
	cropped := NewDrawFromImgColorModel(img.ColorModel(), rect)
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)
	return cropped
}
