package filextract

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
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

func (fp *FileProcessorWorker) Run() {
	for filename := range fp.FilenameStream {
		if err := fp.processFile(filename); err != nil {
			log.Fatal(err)
		}

		// After finishing file processing, start the post analysis
		if err := outdirwriter.MoveFirstFileToCoverFolder(filepath.Join(fp.OutputFolder, filename.BaseName)); err != nil {
			log.Fatal(err)
		}
	}
}

func (fp *FileProcessorWorker) processFile(file FileInfo) error {
	extractDir := filepath.Join(fp.OutputFolder, file.BaseName)

	// Create the directory for the extracted files
	if err := outdirwriter.CreateOutDir(extractDir, outdirwriter.CoverDirSuffix); err != nil {
		return err
	}

	filePointer, err := os.OpenFile(file.CompleteName, os.O_RDONLY, 0755)
	if err != nil {
		log.Printf("Failed to open %s: %v", file.CompleteName, err)
		return err
	}

	var (
		extractor  cbxr.Extractor
		fileWriter = outdirwriter.NewWriterHandle(extractDir)
	)

	if extractor, err = fp.newExtractor(file, filePointer); err != nil {
		log.Printf("Failed to create extractor: %v", err)
		return err
	}

	for fileName, fileResult := range extractor.FileSeq() {
		if fileResult.Error != nil {
			return fileResult.Error
		}

		fileName = cbxr.FileName(filepath.Base(string(fileName)))
		if strings.HasPrefix(strings.ToLower(string(fileName)), "cred") &&
			len(fileName) >= len("000.jpeg") {
			continue
		}

		var decodedImg image.Image
		if decodedImg, _, err = image.Decode(bytes.NewReader(fileResult.Data)); err != nil {
			return err
		}

		finalImg, processErr := fp.imgPipeline.Process(decodedImg)
		if processErr = fileWriter.Handler(string(fileName), func(writer io.Writer) error {
			return jpeg.Encode(writer, finalImg, nil)
		}); processErr != nil {
			return processErr
		}
	}

	err = fileWriter.OnFinish()
	return err
}

func (fp *FileProcessorWorker) newExtractor(
	file FileInfo,
	filePointer *os.File,
) (extractor cbxr.Extractor, err error) {
	switch strings.ToLower(filepath.Ext(file.CompleteName)) {
	case ".pdf":
		return cbxr.NewPDFExtractor(filePointer)
	case ".zip":
		var stat os.FileInfo
		if stat, err = filePointer.Stat(); err != nil {
			return nil, err
		}
		return cbxr.NewCBZExtractor(filePointer, stat.Size())
	default:
		return cbxr.NewMultiZipRarExtractor(file.CompleteName, filePointer)
	}
}
