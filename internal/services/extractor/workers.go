package extractor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type (
	FileInfo struct {
		CompleteName string
		BaseName     string
	}
	FileProcessorWorker struct {
		OutputFolder   string
		FilenameStream chan FileInfo
	}
)

func (fp *FileProcessorWorker) Run() {
	for filename := range fp.FilenameStream {
		if err := fp.processFile(filename); err != nil {
			log.Fatal(err)
		}

		// After finishing file processing, start the post analysis
		if err := moveFirstFileToCoverFolder(filepath.Join(fp.OutputFolder, filename.BaseName)); err != nil {
			log.Fatal(err)
		}
	}
}

func (fp *FileProcessorWorker) processFile(file FileInfo) error {
	extractDir := filepath.Join(fp.OutputFolder, file.BaseName)
	cbzFile := file.CompleteName

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
		coverDirectoryName: filepath.Join(extractDir, coverDirSuffix),
		folderCounter:      &folderInfoCounter{},
	}
	if err = format.Extract(ctx, fileReader, getAllNames(cbzFile), fileWriter.handler); err != nil {
		log.Printf("Failed to extract %s: %v", cbzFile, err)
		return err
	}

	if fileWriter.folderCounter.onRoot > 0 && fileWriter.folderCounter.onSubDir > 0 {
		if fileWriter.folderCounter.cover > 0 {
			err = fmt.Errorf("only one of onRoot and onSubDir may be specified")
			return err
		}
		// Move all content-main to work as _0Cover
		err = errors.Join(
			os.Remove(fileWriter.coverDirectoryName),
			os.Rename(fileWriter.defaultDir(), fileWriter.coverDirectoryName),
		)
	}
	return err
}

func (fp *FileProcessorWorker) createOutDir(extractDir string, suffix string) error {
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	// Create a covers output directory
	if err := os.MkdirAll(filepath.Join(extractDir, suffix), 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	return nil
}
