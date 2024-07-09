package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	resourcesFolder = "./resources"
)

func main() {
	var targetFolder string
	flag.StringVar(&targetFolder, "src", "", "Target folder where files are stored")
	flag.Parse()

	if targetFolder == "" {
		log.Fatal("Target folder is required")
	}

	var (
		inputFolder  = resourcesFolder + "/input_books/" + targetFolder
		outputFolder = resourcesFolder + "/extracted/" + targetFolder
	)

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
		sendChannel = make(chan FileInfo)
	)
	// Create worker pool
	for range 5 {
		wg.Add(1)
		go func() {
			fp := fileProcessorWorker{
				outputFolder:   outputFolder,
				filenameStream: sendChannel,
			}
			defer wg.Done()
			fp.run()
		}()
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".cbz") {
			cbzFile := filepath.Join(inputFolder, file.Name())
			baseName := strings.TrimSuffix(file.Name(), ".cbz")
			sendChannel <- FileInfo{
				baseName:     baseName,
				completeName: cbzFile,
			}
		}
	}
	close(sendChannel)

	wg.Wait()
	log.Printf("Sent %d files", len(files))
}
