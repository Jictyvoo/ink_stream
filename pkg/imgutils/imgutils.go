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
	return min(MaxPixelValue, uint8(max(MinPixelValue, value)))
}

type (
	ColorConverter interface {
		Convert(c color.Color) color.Color
	}
	DrawImageFactory interface {
		CreateDrawImage(colorModel color.Model, bounds image.Rectangle) draw.Image
	}
)
