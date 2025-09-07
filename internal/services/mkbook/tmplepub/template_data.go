package tmplepub

// ImageData holds the values used to fill the EPUB image HTML template.
// Keep fields exported so html/template can access them.
// Only ImageSrc is strictly required for minimal rendering; others are optional.
type ImageData struct {
	BodyStyle      string
	TopMargin      int
	ImageWidth     int
	ImageHeight    int
	ImageSrc       string
	ViewportWidth  int
	ViewportHeight int
	PanelLinks     []PanelLink
	PanelImages    []PanelImage
}

type PanelLink struct {
	ID       string
	TargetID string
	Ordinal  int
}

type PanelImage struct {
	ID    string
	Style string
}
