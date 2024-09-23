package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/Jictyvoo/ink_stream/internal/services/extractor"
)

func main() {
	var inputFolder string
	flag.StringVar(&inputFolder, "src", "", "Target folder where files are stored")
	flag.Parse()

	if inputFolder == "" {
		log.Fatal("Target folder is required")
	}

	var (
		lastFolderName = filepath.Base(inputFolder)
		rootDir        = filepath.Dir(strings.TrimSuffix(strings.Trim(inputFolder, "/"), lastFolderName))
	)
	rootDir = strings.TrimSuffix(strings.Trim(rootDir, "/"), filepath.Base(rootDir))

	var outputFolder = filepath.Join(rootDir, "extracted", lastFolderName)
	fmt.Printf("Using target folder %s\n", inputFolder)
	fmt.Printf("Using output folder %s\n", outputFolder)
	// Ensure the output folder exists
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		log.Fatalf("Failed to create output folder: %v", err)
	}

	// Find all .cbz files in the input folder
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		log.Fatalf("Failed to read input folder: %v", err)
	}

	var (
		wg          sync.WaitGroup
		sendChannel = make(chan extractor.FileInfo)
	)
	// Create worker pool
	for range 5 {
		wg.Add(1)
		go func() {
			fp := extractor.FileProcessorWorker{
				OutputFolder:   outputFolder,
				FilenameStream: sendChannel,
			}
			defer wg.Done()
			fp.Run()
		}()
	}

	allowedFormats := []string{".cbz", ".cbr", ".zip", ".rar"}
	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if slices.Contains(allowedFormats, fileExt) {
			fileAbsolutePath := filepath.Join(inputFolder, file.Name())
			baseName := strings.TrimSuffix(file.Name(), fileExt)
			sendChannel <- extractor.FileInfo{
				BaseName:     baseName,
				CompleteName: fileAbsolutePath,
			}
		}
	}
	close(sendChannel)

	wg.Wait()
	log.Printf("Sent %d files", len(files))
}
