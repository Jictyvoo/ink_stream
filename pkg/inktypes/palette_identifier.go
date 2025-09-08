package inktypes

import (
	"bytes"
	"encoding/hex"
	"image/color"
	"strconv"
)

type PaletteIdentifier []byte

func NewPaletteIdentifier(fromColor color.Palette) PaletteIdentifier {
	shiftPixel := func(elem ...*uint32) {
		for _, pixelPtr := range elem {
			if pixelPtr != nil {
				*pixelPtr = *pixelPtr >> 8
			}
		}
	}

	var buffer bytes.Buffer
	for index, paletteColor := range fromColor {
		r, g, b, a := paletteColor.RGBA()
		shiftPixel(&r, &g, &b, &a)
		color64 := uint64(r) + uint64(g)<<8 + uint64(b)<<16 + uint64(a)<<24
		if index > 0 {
			buffer.WriteRune('-')
		}
		buffer.WriteString(strconv.FormatUint(color64, 16))
	}

	bytesArr := make([]byte, hex.EncodedLen(buffer.Len()))
	_ = hex.Encode(bytesArr, buffer.Bytes())
	return bytesArr
}

func (pi PaletteIdentifier) Hex() string {
	return string(pi)
}
