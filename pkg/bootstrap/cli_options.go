package bootstrap

import (
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

type CropStyle uint8

const (
	CropBasic CropStyle = iota
	CropNormal
	CropAggressive
)

type OutputFormat string

const (
	FormatEpub   OutputFormat = "epub"
	FormatMobi   OutputFormat = "mobi"
	FormatFolder OutputFormat = "folder"
)

type ReadDirection string

type ImageFormat string

const (
	ImageJPEG = ImageFormat(inktypes.FormatJPEG)
	ImagePNG  = ImageFormat(inktypes.FormatPNG)
)

type Options struct {
	SourceFolder  string
	OutputFolder  string
	TargetDevice  deviceprof.DeviceType
	CropLevel     CropStyle
	OutputFormat  OutputFormat
	ReadDirection ReadDirection
	RotateImage   bool
	StretchImage  bool
	AddMargins    bool
	ColoredPages  bool
	ImageFormat   ImageFormat
	ImageQuality  uint8
}

func (opts Options) AllowStretch() bool {
	if opts.AddMargins {
		return true
	}

	return opts.StretchImage
}
