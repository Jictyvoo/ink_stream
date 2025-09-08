package inktypes

type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
	FormatBMP  ImageFormat = "bmp"
	FormatTIFF ImageFormat = "tiff"
	FormatWEBP ImageFormat = "webp"
)

type ImageMetadata struct {
	ImageDimensions
	Format  ImageFormat
	Palette PaletteIdentifier
	DPI     int
}

func NewImageMetadata(width, height int) ImageMetadata {
	return ImageMetadata{
		ImageDimensions: ImageDimensions{
			Width:  uint16(width),
			Height: uint16(height),
		},
	}
}
