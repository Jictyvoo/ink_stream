package main

import (
	"encoding/json"
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
	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
	"github.com/Jictyvoo/ink_stream/internal/utils"
)

func main() {
	var cliArgs Options
	parseArgs(&cliArgs)

	{
		jsonBytes, _ := json.Marshal(cliArgs)
		var argsMap map[string]any
		_ = json.Unmarshal(jsonBytes, &argsMap)
		slog.Info("Using cli options", slog.Any("options", argsMap))
	} // Ensure the output folder exists
	if err := os.MkdirAll(cliArgs.OutputFolder, 755); err != nil {
		log.Fatalf("Failed to create output folder: %v", err)
	}

	var (
		wg          sync.WaitGroup
		sendChannel = make(chan filextract.FileInfo)
	)

	imgPipeline, err := BuildPipeline(cliArgs)
	if err != nil {
		slog.Error("Failed to build pipeline", slog.String("error", err.Error()))
		return
	}
	// Create worker pool
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			fp := filextract.NewFileProcessorWorker(
				sendChannel,
				cliArgs.OutputFolder,
				func(outputDir string) (filextract.FileOutputWriter, error) {
					fileWriter := outdirwriter.NewWriterHandle(outputDir)
					return imgprocessor.NewMultiThreadImageProcessor(fileWriter, imgPipeline), nil
				},
			)
			defer wg.Done()
			_ = fp.Run()
		}()
	}

	filenameList := utils.ListAllFiles(cliArgs.SourceFolder)
	allowedFormats := cbxr.SupportedFileExtensions()
	for _, fileAbsolutePath := range filenameList {
		fileExt := strings.ToLower(filepath.Ext(fileAbsolutePath))
		if slices.Contains(allowedFormats, fileExt) {
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
