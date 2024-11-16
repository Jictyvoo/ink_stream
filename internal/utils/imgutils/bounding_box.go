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

// CropBox calculates a bounding box for the given image by analyzing transparent pixels or minimum color values.
func CropBox(img image.Image, colorConverter ColorConverter, opts BoxOptions) image.Rectangle {
	// Use a default color converter if none is provided.
	if colorConverter == nil {
		colorConverter = color.GrayModel
	}

	// Initialize bounding box and image dimensions.
	var (
		bbox                               = img.Bounds()
		width                              = bbox.Dx()
		height                             = bbox.Dy()
		transparentAnalysis, whiteAnalysis struct {
			row, column []uint64
		}
		whiteValue uint64 // Keeps track of the most prominent "white" value.
	)

	// Allocate slices for analyzing transparency or minimum color if needed.
	if opts.Is(BoxEliminateTransparent) {
		transparentAnalysis.row = make([]uint64, height)
		transparentAnalysis.column = make([]uint64, width)
	}
	if opts.Is(BoxEliminateMinimumColor) {
		whiteAnalysis.row = make([]uint64, height)
		whiteAnalysis.column = make([]uint64, width)
	}

	// Iterate over all pixels in the image using a custom Iterator function.
	for x, y := range Iterator(img) {
		pixel := img.At(x, y)
		_, _, _, a := pixel.RGBA() // Extract alpha value to check transparency.

		// Update transparency analysis if enabled.
		if opts.Is(BoxEliminateTransparent) && a == 0 {
			transparentAnalysis.column[x]++
			transparentAnalysis.row[y]++
		}

		// Update minimum color analysis if enabled.
		if opts.Is(BoxEliminateMinimumColor) {
			convertedPixel := colorConverter.Convert(pixel)
			r, g, b, _ := convertedPixel.RGBA()
			pixelValue := (uint64(r) << 8) | uint64(g) | uint64(b>>8)

			// Track rows and columns matching the current "white" value.
			if whiteValue == 0 || pixelValue >= whiteValue {
				whiteAnalysis.column[x]++
				whiteAnalysis.row[y]++
			}
			whiteValue = max(whiteValue, pixelValue) // Update the maximum white value.
		}
	}

	// Adjust the bounding box based on analysis slices.
	newBox := bbox
	if opts.Is(BoxEliminateMinimumColor) {
		newBox = cutBoxBasedOn(newBox, whiteAnalysis)
	}
	if opts.Is(BoxEliminateTransparent) {
		newBox = cutBoxBasedOn(newBox, transparentAnalysis)
	}
	return newBox
}

// cutBoxBasedOn refines the bounding box by removing rows and columns
// that match the background criteria defined in analysisSlices.
func cutBoxBasedOn(
	originalBox image.Rectangle,
	analysisSlices struct{ row, column []uint64 },
) image.Rectangle {
	width, height := originalBox.Dx(), originalBox.Dy()
	newValues := [2]image.Point{originalBox.Min, originalBox.Max}

	// Helper function to adjust bounding box based on analysis.
	valueChanger := func(valueSlice []uint64, index int, maxValue int, minMaxSel *uint8, destination *int, modifier *int) int {
		workOnMinimum := *minMaxSel == 0
		switch valueSlice[index] {
		case uint64(maxValue):
			*destination += *modifier
		default: // No match, increment minMaxSel to start processing the opposite edge.
			*minMaxSel = min(1, *minMaxSel+1)
			if workOnMinimum {
				*modifier = -1
				index = maxValue
			}
		}

		return index + *modifier
	}

	var (
		valIndexes   struct{ x, y uint8 } // Tracks which edge (min or max) is being processed.
		loopModifier = image.Point{
			X: 1, Y: 1,
		} // Determines scan direction (forward or backward).
		colIndex, rowIndex int // Current column and row indices.
	)

	for range max(width, height) {
		if colIndex < width {
			colIndex = valueChanger(
				analysisSlices.column, colIndex, width,
				&valIndexes.x, &newValues[valIndexes.x].X, &loopModifier.X,
			)
		}
		if rowIndex < height {
			rowIndex = valueChanger(
				analysisSlices.row, rowIndex, height,
				&valIndexes.y, &newValues[valIndexes.y].Y, &loopModifier.Y,
			)
		}
	}

	return image.Rect(newValues[0].X, newValues[0].Y, newValues[1].X, newValues[1].Y)
}
