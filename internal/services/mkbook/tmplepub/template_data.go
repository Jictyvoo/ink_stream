package tmplepub

// ImageData holds the values used to fill the EPUB image HTML template.
// Keep fields exported so html/template can access them.
// Only ImageSrc is strictly required for minimal rendering; others are optional.
type ImageData struct {
	TopMargin      int
	ImageWidth     int
	ImageHeight    int
	ImageSrc       string
	ViewportWidth  int
	ViewportHeight int

	// BaseID is derived from ImageSrc and used to generate unique
	// IDs for tap targets and their corresponding panel elements.
	BaseID string

	// Panel definitions
	PanelImages []PanelImage
}

type PanelImage struct {
	Class   string
	Ordinal int
}
