package extractor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/internal/services/extractor/cbxr"
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

	// Create the directory for the extracted files
	if err := fp.createOutDir(extractDir, coverDirSuffix); err != nil {
		return err
	}

	filePointer, err := os.OpenFile(file.CompleteName, os.O_RDONLY, 0755)
	if err != nil {
		log.Printf("Failed to open %s: %v", file.CompleteName, err)
		return err
	}

	var (
		extractor  cbxr.Extractor
		fileWriter = WriterHandle{
			outputDirectory:    extractDir,
			coverDirectoryName: filepath.Join(extractDir, coverDirSuffix),
			folderCounter:      &folderInfoCounter{},
		}
	)
	extractor, err = cbxr.NewMultiZipRarExtractor(file.CompleteName, filePointer)
	if err != nil {
		log.Printf("Failed to create extractor: %v", err)
		return err
	}

	for fileName, fileResult := range extractor.FileSeq() {
		if fileResult.Error != nil {
			return fileResult.Error
		}

		fileName = cbxr.FileName(filepath.Base(string(fileName)))
		if strings.HasPrefix(strings.ToLower(string(fileName)), "cred") && len(fileName) >= len("000.jpeg") {
			continue
		}
		if err = fileWriter.handler(string(fileName), fileResult.Data); err != nil {
			return err
		}
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
