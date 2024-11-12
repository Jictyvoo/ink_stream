package filextract

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/internal/deviceprof"
	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/internal/services/filextract/cbxr"
	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
)

type (
	FileInfo struct {
		CompleteName string
		BaseName     string
	}
	FileProcessorWorker struct {
		OutputFolder   string
		FilenameStream chan FileInfo
		imgPipeline    imageparser.ImagePipeline
	}
)

func NewFileProcessorWorker(
	filenameStream chan FileInfo,
	outputFolder string,
	targetResolution deviceprof.Resolution,
) *FileProcessorWorker {
	return &FileProcessorWorker{
		FilenameStream: filenameStream,
		OutputFolder:   outputFolder,
		imgPipeline: imageparser.NewImagePipeline(
			imageparser.NewStepAutoContrast(0, 0),
			imageparser.NewStepRescale(targetResolution),
			imageparser.NewStepGrayScale(),
		),
	}
}

func (fp *FileProcessorWorker) Run() error {
	for filename := range fp.FilenameStream {
		if err := fp.processFile(filename); err != nil {
			return err
		}

		// After finishing file processing, start the post analysis
		if err := outdirwriter.MoveFirstFileToCoverFolder(filepath.Join(fp.OutputFolder, filename.BaseName)); err != nil {
			return err
		}
	}

	return nil
}

func (fp *FileProcessorWorker) processFile(file FileInfo) (resultErr error) {
	extractDir := filepath.Join(fp.OutputFolder, file.BaseName)

	// Create the directory for the extracted files
	if err := outdirwriter.CreateOutDir(extractDir, outdirwriter.CoverDirSuffix); err != nil {
		return err
	}

	filePointer, err := os.OpenFile(file.CompleteName, os.O_RDONLY, 0755)
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
		extractor       cbxr.Extractor
		multiThreadProc = NewMultiThreadImageProcessor(extractDir, fp.imgPipeline)
	)
	defer multiThreadProc.Close()

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

		multiThreadProc.Process(string(fileName), fileResult.Data)
		totalSent++
	}

	err = multiThreadProc.Shutdown()
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
	case ".pdf":
		return cbxr.NewPDFExtractor(filePointer)
	case ".zip":
		return cbxr.NewCBZExtractor(filePointer)
	default:
		return cbxr.NewMultiZipRarExtractor(file.CompleteName, filePointer)
	}
}
