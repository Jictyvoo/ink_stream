package inktypes

import "strings"

type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
	FormatBMP  ImageFormat = "bmp"
	FormatTIFF ImageFormat = "tiff"
	FormatWEBP ImageFormat = "webp"
)

type (
	ImageEncodingOptions struct {
		Quality uint8
		Format  ImageFormat
	}
	ImageMetadata struct {
		ImageDimensions
		ImageEncodingOptions
		Palette PaletteIdentifier
		DPI     int
	}
)

func NewImageEncodingOptions(quality uint8, format ImageFormat) ImageEncodingOptions {
	quality = min(max(quality, 60), 100)
	format = ImageFormat(strings.ToLower(string(format)))
	if format == "" {
		format = FormatJPEG
	}
	return ImageEncodingOptions{Quality: quality, Format: format}
}

func NewImageMetadata(width, height int) ImageMetadata {
	return ImageMetadata{
		ImageDimensions: ImageDimensions{
			Width:  uint16(width),
			Height: uint16(height),
		},
	}
}

func (ieo ImageEncodingOptions) FileExtension() string {
	return "." + string(ieo.Format)
}
