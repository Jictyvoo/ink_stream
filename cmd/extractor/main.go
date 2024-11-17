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

	"github.com/Jictyvoo/ink_stream/internal/services/filextract"
	"github.com/Jictyvoo/ink_stream/internal/utils"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
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
		// baseDir := filepath.Base(rootDir)
		// rootDir = strings.TrimSuffix(strings.TrimSuffix(rootDir, "/"), baseDir)
		outputFolder = filepath.Join(rootDir, "extracted", lastFolderName)
	}
	fmt.Printf("Using target folder %s\n", inputFolder)
	fmt.Printf("Using output folder %s\n", outputFolder)
	// Ensure the output folder exists
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		log.Fatalf("Failed to create output folder: %v", err)
	}

	var (
		wg          sync.WaitGroup
		sendChannel = make(chan filextract.FileInfo)
	)

	deviceProfile, _ := deviceprof.Profile(deviceprof.DeviceKindle11)
	// Create worker pool
	for range 5 {
		wg.Add(1)
		go func() {
			fp := filextract.NewFileProcessorWorker(
				sendChannel, outputFolder, deviceProfile,
			)
			defer wg.Done()
			_ = fp.Run()
		}()
	}

	filenameList := utils.ListAllFiles(inputFolder)
	allowedFormats := []string{".cbz", ".cbr", ".zip", ".rar", ".pdf"}
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
