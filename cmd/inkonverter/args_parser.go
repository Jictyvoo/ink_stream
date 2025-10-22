package main

import (
	"errors"
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/pkg/bootstrap"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func parseArgs(cliArgs *bootstrap.Options) {
	flag.StringVar(&cliArgs.SourceFolder, "src", "", "Target folder where files are stored")
	flag.StringVar(&cliArgs.OutputFolder, "out", "", "Output folder where files will be saved")
	flag.BoolVar(&cliArgs.RotateImage, "rotate", false, "Rotate image files")
	flag.BoolVar(&cliArgs.ColoredPages, "colored", false, "Colored pages")
	flag.BoolVar(&cliArgs.AddMargins, "margins", false, "Add margin on image")
	flag.BoolVar(&cliArgs.StretchImage, "stretch", true, "Stretch image files")
	cropLevel := flag.Uint("crop-level", uint(bootstrap.CropBasic), "Crop image level")

	var (
		targetDevice  string
		outFormat     string
		readDirection string
		imgOutFormat  string
		imgOutQuality uint
	)
	flag.StringVar(&targetDevice, "profile", "", "Target device name")
	flag.StringVar(&outFormat, "format", string(bootstrap.FormatEpub), "Output format")
	flag.StringVar(&imgOutFormat, "img-format", string(bootstrap.ImageJPEG), "Image output format")
	flag.UintVar(&imgOutQuality, "img-quality", 85, "Image output quality")
	flag.StringVar(
		&readDirection, "read-direction",
		inktypes.ReadLeftToRight.String(), "Read direction used as epub PPD",
	)
	flag.Parse()

	cliArgs.CropLevel = bootstrap.CropBasic
	if cropLevel != nil {
		cliArgs.CropLevel = bootstrap.CropStyle(*cropLevel)
	}
	cliArgs.TargetDevice = deviceprof.DeviceType(targetDevice)
	cliArgs.OutputFormat = bootstrap.OutputFormat(outFormat)
	cliArgs.ImageQuality = uint8(imgOutQuality)
	cliArgs.ImageFormat = bootstrap.ImageFormat(imgOutFormat)
	if cliArgs.ImageFormat == bootstrap.ImagePNG {
		cliArgs.ImageQuality = 100
	}

	if cliArgs.OutputFolder == "" {
		cliArgs.OutputFolder = defaultOutputFolder(cliArgs.SourceFolder)
	}
	cliArgs.ReadDirection = bootstrap.ReadDirection(readDirection)
	cliErr := func(err error) {
		flag.Usage()
		log.Fatal(err)
	}
	if cliArgs.SourceFolder == "" {
		cliErr(errors.New("target folder is required"))
	}
	if cliArgs.TargetDevice == "" {
		cliErr(errors.New("target device is required"))
	}
	if cliArgs.OutputFormat == "" {
		cliErr(errors.New("output format is required"))
	}
}

func defaultOutputFolder(srcDir string) string {
	var (
		lastFolderName = filepath.Base(srcDir)
		rootDir        = filepath.Dir(
			strings.TrimSuffix(strings.TrimSuffix(srcDir, "/"), lastFolderName),
		)
	)

	return filepath.Join(rootDir, "converted", lastFolderName)
}
