package imgutils

import (
	"image"
	"image/color"
)

// DominantColorInRegion computes the average color in a region and returns the nearest color from the palette of that region.
func DominantColorInRegion(
	img image.Image,
	region image.Rectangle,
	optTakeAvg ...bool,
) color.Color {
	var totals struct{ r, g, b, a uint64 }
	var count uint64

	palette := make(color.Palette, 0, 4)
	colorsSet := make(map[color.NRGBA]struct{})
	for x, y := range RegionIterator(region) {
		r, g, b, a := img.At(x, y).RGBA()
		if a == 0 {
			continue // skip fully transparent
		}
		c := color.NRGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
		totals.r += uint64(c.R)
		totals.g += uint64(c.G)
		totals.b += uint64(c.B)
		totals.a += uint64(c.A)
		count++

		if _, ok := colorsSet[c]; !ok {
			palette = append(palette, c)
			colorsSet[c] = struct{}{}
		}
	}

	clear(colorsSet) // Free memory as soon as possible
	if count == 0 || len(palette) == 0 {
		return color.NRGBA{}
	}

	avg := color.NRGBA{
		R: uint8(totals.r / count),
		G: uint8(totals.g / count),
		B: uint8(totals.b / count),
		A: uint8(totals.a / count),
	}
	if len(optTakeAvg) > 0 && optTakeAvg[0] {
		return avg
	}

	// Find nearest color in palette to average
	return palette.Convert(avg)
}

func ImageMarginDominantColor(
	inputImg image.Image, xMargin, yMargin uint32, percentage uint8,
) (dominantColors Margins[color.Color]) {
	bounds := inputImg.Bounds()
	imageRes := struct{ Width, Height int }{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}
	// Get dominant colors per margin
	if yMargin > 0 {
		heightRange := (imageRes.Height * int(percentage)) / 100 // ?% of the height
		if heightRange <= 0 {                                    // Fallback case
			heightRange = imageRes.Height
		}
		dominantColors.Top = DominantColorInRegion(
			inputImg, image.Rect(0, 0, imageRes.Width, heightRange),
		)
		dominantColors.Bottom = DominantColorInRegion(
			inputImg, image.Rect(0, imageRes.Height-heightRange, imageRes.Width, imageRes.Height),
		)
	}
	if xMargin > 0 {
		widthRange := (imageRes.Width * int(percentage)) / 100 // ?% of the width
		if widthRange <= 0 {                                   // Fallback case
			widthRange = imageRes.Width
		}
		dominantColors.Left = DominantColorInRegion(
			inputImg, image.Rect(0, 0, widthRange, imageRes.Height),
		)
		dominantColors.Right = DominantColorInRegion(
			inputImg, image.Rect(imageRes.Width-widthRange, 0, imageRes.Width, imageRes.Height),
		)
	}

	return
}
