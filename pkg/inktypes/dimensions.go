package inktypes

type ImageOrientation uint8

const (
	OrientationPortrait ImageOrientation = iota
	OrientationLandscape
)

type ImageDimensions struct {
	Width, Height uint16
}

func NewOrientation(dimensions ImageDimensions) ImageOrientation {
	imgOrientation := OrientationPortrait
	if dimensions.Width > dimensions.Height {
		imgOrientation = OrientationLandscape
	}

	return imgOrientation
}

func (dim ImageDimensions) Orientation() ImageOrientation {
	return NewOrientation(dim)
}
