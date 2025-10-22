package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/Jictyvoo/ink_stream/internal/services/filextract"
	"github.com/Jictyvoo/ink_stream/internal/services/filextract/cbxr"
	"github.com/Jictyvoo/ink_stream/internal/utils"
	"github.com/Jictyvoo/ink_stream/pkg/bootstrap"
)

func main() {
	var (
		inputFolder  string
		outputFolder string
	)
	flag.StringVar(&inputFolder, "src", "", "Target folder where files are stored")
	flag.StringVar(&outputFolder, "out", "", "Output folder where files will be saved")
	flag.Parse()

	if inputFolder == "" {
		log.Fatal("Target folder is required")
	}

	var (
		lastFolderName = filepath.Base(inputFolder)
		rootDir        = filepath.Dir(
			strings.TrimSuffix(strings.TrimSuffix(inputFolder, "/"), lastFolderName),
		)
	)

	if outputFolder == "" {
		outputFolder = filepath.Join(rootDir, "extracted", lastFolderName)
	}
	fmt.Printf("Using target folder %s\n", inputFolder)
	fmt.Printf("Using output folder %s\n", outputFolder)
	// Ensure the output folder exists
	if err := os.MkdirAll(outputFolder, 0o755); err != nil {
		log.Fatalf("Failed to create output folder: %v", err)
	}

	var (
		wg          sync.WaitGroup
		sendChannel = make(chan filextract.FileInfo)
	)

	// Create worker pool
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			fp := filextract.NewFileProcessorWorker(
				sendChannel, outputFolder,
				func(outputDir string) (filextract.FileOutputWriter, error) {
					return bootstrap.NewFileWriterWrapper(outputDir)
				},
			)
			defer wg.Done()
			_ = fp.Run()
		}()
	}

	filenameList := utils.ListAllFiles(inputFolder)
	allowedFormats := cbxr.SupportedFileExtensions()
	var sentCount uint16
	for _, fileAbsolutePath := range filenameList {
		fileExt := strings.ToLower(filepath.Ext(fileAbsolutePath))
		if slices.Contains(allowedFormats, fileExt) {
			baseName := strings.TrimSuffix(filepath.Base(fileAbsolutePath), fileExt)
			sendChannel <- filextract.FileInfo{
				BaseName:     baseName,
				CompleteName: fileAbsolutePath,
			}
			sentCount++
		}
	}
	close(sendChannel)

	wg.Wait()
	log.Printf("Sent %d files", sentCount)
}
