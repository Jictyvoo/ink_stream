package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

type (
	FileInfo struct {
		completeName string
		baseName     string
	}
	fileProcessorWorker struct {
		outputFolder   string
		filenameStream chan FileInfo
	}
)

func (fp *fileProcessorWorker) run() {
	for filename := range fp.filenameStream {
		if err := fp.processFile(filename); err != nil {
			log.Fatal(err)
		}
	}
}

func (fp *fileProcessorWorker) processFile(file FileInfo) error {
	const coverDirSuffix = "_0Cover"
	extractDir := filepath.Join(outputFolder, file.baseName)
	cbzFile := file.completeName

	// Create the directory for the extracted files
	if err := fp.createOutDir(extractDir, coverDirSuffix); err != nil {
		return err
	}

	filePointer, err := os.OpenFile(cbzFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Failed to open %s: %v", cbzFile, err)
		return err
	}

	fileReader, format, err := checkFileFormat(cbzFile, filePointer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 9000*time.Second)
	defer cancel()

	fileWriter := WriterHandle{
		outputDirectory:    extractDir,
		coverDirectoryName: extractDir + "/" + coverDirSuffix,
	}
	if err = format.Extract(ctx, fileReader, getAllNames(cbzFile), fileWriter.handler); err != nil {
		log.Printf("Failed to extract %s: %v", cbzFile, err)
		return err
	}

	return nil
}

func (fp *fileProcessorWorker) createOutDir(extractDir string, suffix string) error {
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	// Create a covers output directory
	if err := os.MkdirAll(extractDir+"/"+suffix, 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	return nil
}
