package main

import (
	"errors"
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func parseArgs(cliArgs *Options) {
	flag.StringVar(&cliArgs.SourceFolder, "src", "", "Target folder where files are stored")
	flag.StringVar(&cliArgs.OutputFolder, "out", "", "Output folder where files will be saved")
	flag.BoolVar(&cliArgs.RotateImage, "rotate", false, "Rotate image files")
	flag.BoolVar(&cliArgs.ColoredPages, "colored", false, "Colored pages")
	cliArgs.AddMargins = flag.Bool("margins", false, "Add margin on image")
	cliArgs.StretchImage = flag.Bool("stretch", false, "Stretch image files")
	cropLevel := flag.Uint("crop-level", uint(CropBasic), "Crop image level")

	var (
		targetDevice  string
		outFormat     string
		readDirection string
	)
	flag.StringVar(&targetDevice, "profile", "", "Target device name")
	flag.StringVar(&outFormat, "format", string(FormatEpub), "Output format")
	flag.StringVar(
		&readDirection, "read-direction",
		inktypes.ReadLeftToRight.String(), "Read direction used as epub PPD",
	)
	flag.Parse()

	cliArgs.CropLevel = CropBasic
	if cropLevel != nil {
		cliArgs.CropLevel = CropStyle(*cropLevel)
	}
	cliArgs.TargetDevice = deviceprof.DeviceType(targetDevice)
	cliArgs.OutputFormat = OutputFormat(outFormat)

	if cliArgs.OutputFolder == "" {
		cliArgs.OutputFolder = defaultOutputFolder(cliArgs.SourceFolder)
	}
	cliArgs.ReadDirection = ReadDirection(readDirection)
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
