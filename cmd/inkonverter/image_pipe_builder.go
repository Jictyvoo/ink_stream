package main

import (
	"errors"
	"image/color"
	"slices"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/internal/imageparser/imgpipesteps"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func BuildPipeline(opts Options) (imageparser.ImagePipeline, error) {
	targetProfile, ok := deviceprof.Profile(opts.TargetDevice)
	if !ok {
		return imageparser.ImagePipeline{}, errors.New("target device not found")
	}

	readDirection := inktypes.NewReadDirection(string(opts.ReadDirection))
	autocropPalette := genPalette(opts.CropLevel, targetProfile.Palette)
	imgSteps := append(
		make([]imageparser.PipeStep, 0, 6),
		imgpipesteps.NewStepAutoCrop(autocropPalette),
		imgpipesteps.NewStepMarginWrap(targetProfile.Resolution),
		imgpipesteps.NewStepCropOrRotate(
			opts.RotateImage, color.Palette(targetProfile.Palette),
			readDirection, targetProfile.Resolution.Orientation(),
		),
		imgpipesteps.NewStepRescale(targetProfile.Resolution, opts.AllowStretch()),
		imgpipesteps.NewStepAutoContrast(0, 0),
	)

	if opts.RotateImage { // Only include the margin first if the image should not rotate
		imgSteps[1], imgSteps[2] = imgSteps[2], imgSteps[1]
	}

	if !opts.AddMargins {
		imgSteps = slices.DeleteFunc(imgSteps, func(step imageparser.PipeStep) bool {
			_, isType := step.(*imgpipesteps.StepMarginWrapImage)
			return isType
		})
	}
	if !opts.ColoredPages {
		imgSteps = append([]imageparser.PipeStep{imgpipesteps.NewStepGrayScale()}, imgSteps...)
	}

	builtPipe := imageparser.NewImagePipeline(
		color.Palette(targetProfile.Palette), imgSteps...,
	)

	return builtPipe, nil
}

func genPalette(level CropStyle, palette deviceprof.PaletteType) color.Palette {
	switch level {
	case CropBasic:
		return color.Palette(palette)
	case CropNormal:
		return color.Palette{
			color.Black, color.White,
			color.Gray16{Y: 0xaaaa}, // Light gray (~66%)
			color.Gray16{Y: 0x5555}, // Dark gray (~33%)
		}
	case CropAggressive:
		return color.Palette{
			color.Black, color.White,
		}
	}

	return color.Palette(palette)
}
