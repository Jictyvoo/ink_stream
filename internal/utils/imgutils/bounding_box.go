package imgutils

import (
	"image"
	"image/color"
)

type BoxOptions uint8

func (bo BoxOptions) Is(opt BoxOptions) bool {
	return bo&opt == opt
}

const (
	BoxEliminateTransparent BoxOptions = 1 << iota
	BoxEliminateMinimumColor
)

func CropBox(img image.Image, colorConverter ColorConverter, opts BoxOptions) image.Rectangle {
	if colorConverter == nil {
		colorConverter = color.GrayModel
	}

	var (
		bbox                               = img.Bounds()
		width                              = bbox.Dx()
		height                             = bbox.Dy()
		transparentAnalysis, whiteAnalysis struct {
			row, column []uint64
		}
		whiteValue uint8
	)

	if opts.Is(BoxEliminateTransparent) {
		transparentAnalysis.row = make([]uint64, height)
		transparentAnalysis.column = make([]uint64, width)
	}
	if opts.Is(BoxEliminateMinimumColor) {
		whiteAnalysis.row = make([]uint64, height)
		whiteAnalysis.column = make([]uint64, width)
	}
	// Define initial bounding box by scanning non-background pixels (this part depends on image content)
	for x := range width {
		for y := range height {
			pixel := img.At(x, y)
			_, _, _, a := pixel.RGBA()
			if opts.Is(BoxEliminateTransparent) && a == 0 {
				transparentAnalysis.column[x]++
				transparentAnalysis.row[y]++
			}
			if opts.Is(BoxEliminateMinimumColor) {
				convertedPixel := colorConverter.Convert(pixel)
				r, g, b, _ := convertedPixel.RGBA()
				pixelValue := uint8((r>>8)+(g>>8)+(b>>8)) / 3
				if whiteValue == 0 || pixelValue == whiteValue {
					whiteAnalysis.column[x]++
					whiteAnalysis.row[y]++
				}
				whiteValue = max(whiteValue, pixelValue)
			}
		}
	}

	newBox := bbox
	if opts.Is(BoxEliminateMinimumColor) {
		newBox = cutBoxBasedOn(newBox, whiteAnalysis)
	}
	if opts.Is(BoxEliminateTransparent) {
		newBox = cutBoxBasedOn(newBox, transparentAnalysis)
	}
	return newBox
}

func cutBoxBasedOn(
	originalBox image.Rectangle,
	analysisSlices struct{ row, column []uint64 },
) image.Rectangle {
	width, height := originalBox.Dx(), originalBox.Dy()
	newValues := [2]struct{ x, y int }{{}, {width, height}}
	constructors := func(value uint64, maxValue int, index *uint8, destination *int) {
		switch value {
		case uint64(maxValue):
			modifier := 1
			if *index > 0 {
				modifier *= -1
			}
			*destination += modifier
		default:
			*index = min(1, *index+1)
		}
	}

	var valIndexes struct{ x, y uint8 }
	for index := range max(width, height) {
		if index < width {
			constructors(
				analysisSlices.column[index], width,
				&valIndexes.x, &newValues[valIndexes.x].x,
			)
		}
		if index < height {
			constructors(
				analysisSlices.row[index], height,
				&valIndexes.y, &newValues[valIndexes.y].y,
			)
		}
	}

	return image.Rect(newValues[0].x, newValues[0].y, newValues[1].x, newValues[1].y)
}
