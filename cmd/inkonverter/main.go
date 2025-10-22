package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/Jictyvoo/ink_stream/internal/services/filextract"
	"github.com/Jictyvoo/ink_stream/internal/services/filextract/cbxr"
	"github.com/Jictyvoo/ink_stream/internal/services/imgprocessor"
	"github.com/Jictyvoo/ink_stream/internal/services/mkbook"
	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
	"github.com/Jictyvoo/ink_stream/internal/utils"
	"github.com/Jictyvoo/ink_stream/pkg/bootstrap"
	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

func main() {
	var cliArgs bootstrap.Options
	parseArgs(&cliArgs)

	{
		jsonBytes, _ := json.Marshal(cliArgs)
		var argsMap map[string]any
		_ = json.Unmarshal(jsonBytes, &argsMap)
		slog.Info("Using cli options", slog.Any("options", argsMap))
	} // Ensure the output folder exists
	if err := os.MkdirAll(cliArgs.OutputFolder, 0o755); err != nil {
		log.Fatalf("Failed to create output folder: %v", err)
	}

	var (
		wg          sync.WaitGroup
		sendChannel = make(chan filextract.FileInfo)
	)

	imgPipeline, err := bootstrap.BuildPipeline(cliArgs)
	if err != nil {
		slog.Error("Failed to build pipeline", slog.String("error", err.Error()))
		return
	}

	outWriterFactory, newWriterErr := fileWriterGenerator(
		cliArgs.OutputFormat, cliArgs.ReadDirection,
	)
	if newWriterErr != nil {
		slog.Error("Failed to create output writer", slog.String("error", newWriterErr.Error()))
		os.Exit(1)
	}
	// Create worker pool
	for index := range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			fp := filextract.NewFileProcessorWorker(
				sendChannel, cliArgs.OutputFolder,
				func(outputDir string) (filextract.FileOutputWriter, error) {
					fileWriter, constructErr := outWriterFactory(outputDir)
					imageProcessor := imgprocessor.NewMultiThreadImageProcessor(
						imgPipeline,
						fileWriter, inktypes.NewImageEncodingOptions(
							cliArgs.ImageQuality,
							inktypes.ImageFormat(cliArgs.ImageFormat),
						),
					)
					return imageProcessor, constructErr
				},
			)
			defer wg.Done()
			if processErr := fp.Run(); processErr != nil {
				slog.Error(
					fmt.Sprintf("Failed to process file, goroutine #%d finished", index),
					slog.String("error", processErr.Error()),
					slog.Int("remaining_goroutines", runtime.NumGoroutine()),
				)
			}
		}()
	}

	filenameList := utils.ListAllFiles(cliArgs.SourceFolder)
	allowedFormats := cbxr.SupportedFileExtensions()
	filenameList = utils.CollapseFilesByExt(filenameList, imgutils.SupportedImageFormats())
	for _, fileAbsolutePath := range filenameList {
		fileExt := strings.ToLower(filepath.Ext(fileAbsolutePath))
		if fileExt == "" || slices.Contains(allowedFormats, fileExt) {
			baseName := strings.TrimSuffix(filepath.Base(fileAbsolutePath), fileExt)
			sendChannel <- filextract.FileInfo{
				BaseName:     baseName,
				CompleteName: fileAbsolutePath,
			}
		}
	}
	close(sendChannel)

	wg.Wait()
	log.Printf("Sent %d files", len(filenameList))
}

func fileWriterGenerator(
	format bootstrap.OutputFormat, readDirection bootstrap.ReadDirection,
) (func(outputDir string) (imgprocessor.FileWriter, error), error) {
	direction := inktypes.NewReadDirection(string(readDirection))
	if direction == inktypes.ReadUnknown {
		return nil, fmt.Errorf("unknown read direction `%s`", readDirection)
	}
	switch format {
	case bootstrap.FormatFolder:
		return func(outputDir string) (imgprocessor.FileWriter, error) {
			return outdirwriter.NewWriterHandle(outputDir)
		}, nil
	case bootstrap.FormatEpub:
		return func(outputDir string) (imgprocessor.FileWriter, error) {
			return mkbook.NewEpubMounter(outputDir, direction)
		}, nil
	case bootstrap.FormatMobi:
		return nil, errors.New("mobi format not supported yet")
	default:
		return nil, errors.New("unknown output format")
	}
}
