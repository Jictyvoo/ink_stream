package inktypes

import (
	"bytes"
	"encoding/hex"
	"image/color"
	"strconv"
	"strings"
)

func shiftPixel(elem ...*uint32) {
	for _, pixelPtr := range elem {
		if pixelPtr != nil {
			*pixelPtr = *pixelPtr >> 8
		}
	}
}

type PaletteIdentifier []byte

func NewPaletteIdentifier(fromColor color.Palette) PaletteIdentifier {
	var buffer bytes.Buffer
	for index, paletteColor := range fromColor {
		r, g, b, a := paletteColor.RGBA()
		shiftPixel(&r, &g, &b, &a)
		color64 := uint64(r) | uint64(g)<<8 | uint64(b)<<16 | uint64(a)<<24
		if index > 0 {
			buffer.WriteRune('-')
		}
		buffer.WriteString(strconv.FormatUint(color64, 16))
	}

	bytesArr := make([]byte, hex.EncodedLen(buffer.Len()))
	_ = hex.Encode(bytesArr, buffer.Bytes())
	// return buffer.Bytes()
	return bytesArr
}

func (pi PaletteIdentifier) Hex() string {
	return string(pi)
}

// ToPalette reconstructs a color.Palette from the PaletteIdentifier
func (pi PaletteIdentifier) ToPalette() (color.Palette, error) {
	if len(pi) == 0 {
		return color.Palette{}, nil
	}

	// First decode the hex string back to the original buffer
	rawBytes := make([]byte, hex.DecodedLen(len(pi)))
	_, err := hex.Decode(rawBytes, pi)
	if err != nil {
		return nil, err
	}

	// Split on '-'
	parts := strings.Split(string(rawBytes), "-")
	palette := make(color.Palette, 0, len(parts))

	for _, part := range parts {
		// Parse the hex string into a number
		var val uint64
		if val, err = strconv.ParseUint(part, 16, 64); err != nil {
			return nil, err
		}

		// Extract RGBA components
		r := uint8(val & 0xFF)
		g := uint8((val >> 8) & 0xFF)
		b := uint8((val >> 16) & 0xFF)
		a := uint8((val >> 24) & 0xFF)

		palette = append(palette, color.RGBA{R: r, G: g, B: b, A: a})
	}

	return palette, nil
}
