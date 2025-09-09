package main

import "github.com/Jictyvoo/ink_stream/pkg/deviceprof"

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

type Options struct {
	SourceFolder  string
	OutputFolder  string
	TargetDevice  deviceprof.DeviceType
	CropLevel     CropStyle
	OutputFormat  OutputFormat
	ReadDirection ReadDirection
	RotateImage   bool
	StretchImage  *bool
	AddMargins    *bool
	ColoredPages  bool
}

func (opts Options) AllowStretch() bool {
	stretchImg := opts.StretchImage != nil && *opts.StretchImage
	addMargins := opts.AddMargins != nil && *opts.AddMargins
	if addMargins {
		return false
	}

	return stretchImg
}
