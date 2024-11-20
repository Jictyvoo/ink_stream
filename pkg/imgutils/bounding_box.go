package imgutils

import (
	"image"
	"image/color"
)

type BoxOptions uint8

func (bo BoxOptions) Has(opt BoxOptions) bool {
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
	if opts.Has(BoxEliminateTransparent) {
		transparentAnalysis.row = make([]uint64, height)
		transparentAnalysis.column = make([]uint64, width)
	}
	if opts.Has(BoxEliminateMinimumColor) {
		whiteAnalysis.row = make([]uint64, height)
		whiteAnalysis.column = make([]uint64, width)
	}

	// Iterate over all pixels in the image using a custom Iterator function.
	for x, y := range Iterator(img) {
		pixel := img.At(x, y)
		_, _, _, a := pixel.RGBA() // Extract alpha value to check transparency.

		// Update transparency analysis if enabled.
		if opts.Has(BoxEliminateTransparent) && a == 0 {
			transparentAnalysis.column[x]++
			transparentAnalysis.row[y]++
		}

		// Update minimum color analysis if enabled.
		if opts.Has(BoxEliminateMinimumColor) {
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
	switch {
	case opts.Has(BoxEliminateMinimumColor):
		return cutBoxBasedOn(bbox, whiteAnalysis)
	case opts.Has(BoxEliminateTransparent):
		return cutBoxBasedOn(bbox, transparentAnalysis)
	}

	return bbox
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
			*minMaxSel = min(2, *minMaxSel+1)
			if workOnMinimum {
				*modifier = -1
				index = len(valueSlice)
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
		if colIndex < width && colIndex >= 0 && valIndexes.x <= 1 {
			colIndex = valueChanger(
				analysisSlices.column, colIndex, height,
				&valIndexes.x, &newValues[valIndexes.x].X, &loopModifier.X,
			)
		}
		if rowIndex < height && rowIndex >= 0 && valIndexes.y <= 1 {
			rowIndex = valueChanger(
				analysisSlices.row, rowIndex, width,
				&valIndexes.y, &newValues[valIndexes.y].Y, &loopModifier.Y,
			)
		}
	}

	if newValues[0] == newValues[1] {
		return originalBox
	}

	return image.Rect(newValues[0].X, newValues[0].Y, newValues[1].X, newValues[1].Y)
}

func MarginBox(bounds image.Rectangle, percentage float64) image.Rectangle {
	width, height := bounds.Dx(), bounds.Dy()

	// Calculate minimum and maximum margins
	marginSize := min(int(percentage*float64(width)), int(percentage*float64(height)))

	boundingBox := image.Rect(
		max(0, min(bounds.Min.X-marginSize, bounds.Min.X)),
		max(0, min(bounds.Min.Y-marginSize, bounds.Min.Y)),
		bounds.Max.X+marginSize, bounds.Max.Y+marginSize,
	)

	return boundingBox
}
