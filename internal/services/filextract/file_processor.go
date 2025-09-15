package filextract

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/internal/services/filextract/cbxr"
)

type FileProcessorWorker struct {
	OutputFolder   string
	FilenameStream chan FileInfo
	fileProcessFac FileOutputFactory
}

func NewFileProcessorWorker(
	filenameStream chan FileInfo,
	outputFolder string,
	fileProcessFac FileOutputFactory,
) *FileProcessorWorker {
	return &FileProcessorWorker{
		FilenameStream: filenameStream,
		OutputFolder:   outputFolder,
		fileProcessFac: fileProcessFac,
	}
}

func (fp *FileProcessorWorker) Run() error {
	for filename := range fp.FilenameStream {
		if err := fp.processFile(filename); err != nil {
			return fmt.Errorf("failed to process file `%s`: %w", filename, err)
		}
	}

	return nil
}

func (fp *FileProcessorWorker) processFile(file FileInfo) (resultErr error) {
	extractDir := filepath.Join(fp.OutputFolder, file.BaseName)

	filePointer, err := os.OpenFile(file.CompleteName, os.O_RDONLY, 0o755)
	if err != nil {
		slog.Error(
			"Failed to open input file",
			slog.String("filename", file.CompleteName),
			slog.String("error", err.Error()),
		)
		return err
	}
	defer func(filePointer *os.File) {
		if err = filePointer.Close(); err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	}(filePointer)

	var (
		extractor           cbxr.Extractor
		fileOutputProcessor FileOutputWriter
	)
	if fileOutputProcessor, err = fp.fileProcessFac(extractDir); err != nil {
		return fmt.Errorf("failed to create file output processor: %w", err)
	}
	defer fileOutputProcessor.Close()

	if extractor, err = fp.newExtractor(file, filePointer); err != nil {
		slog.Error("Failed to create extractor", slog.String("error", err.Error()))
		return err
	}

	var totalSent uint64
	for fileName, fileResult := range extractor.FileSeq() {
		if fileResult.Error != nil {
			return fileResult.Error
		}

		fileName = cbxr.FileName(filepath.Base(string(fileName)))
		if strings.HasPrefix(strings.ToLower(string(fileName)), "cred") &&
			len(fileName) >= len("000.jpeg") {
			continue
		}

		fileOutputProcessor.Process(string(fileName), fileResult.Data)
		totalSent++
	}

	err = fileOutputProcessor.Shutdown()
	slog.Info(
		fmt.Sprintf("Sent a total of %d files", totalSent),
		slog.String("inputFile", file.CompleteName),
	)
	return err
}

func (fp *FileProcessorWorker) newExtractor(
	file FileInfo, filePointer *os.File,
) (extractor cbxr.Extractor, err error) {
	switch strings.ToLower(filepath.Ext(file.CompleteName)) {
	case "":
		return cbxr.NewFolderExtractor(filePointer)
	case ".pdf":
		return cbxr.NewPDFExtractor(filePointer)
	case ".zip":
		return cbxr.NewCBZExtractor(filePointer)
	default:
		return cbxr.NewMultiZipRarExtractor(file.CompleteName, filePointer)
	}
}
