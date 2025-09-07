package imgprocessor

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/internal/utils"
)

type (
	fileEntry                 utils.Entry2[string, []byte]
	MultiThreadImageProcessor struct {
		fileWriter    FileWriter
		imgPipeline   imageparser.ImagePipeline
		wg            sync.WaitGroup
		inputChan     chan fileEntry
		numGoroutines uint8
		isFinished    atomic.Bool
	}
)

func NewMultiThreadImageProcessor(
	fileWriter FileWriter,
	imgPipeline imageparser.ImagePipeline,
) *MultiThreadImageProcessor {
	mtip := &MultiThreadImageProcessor{
		fileWriter:    fileWriter,
		imgPipeline:   imgPipeline,
		numGoroutines: 10,
		inputChan:     make(chan fileEntry),
	}

	for range mtip.numGoroutines {
		mtip.wg.Add(1)
		go mtip.workerHandl()
	}

	return mtip
}

func (mtip *MultiThreadImageProcessor) Process(filename string, data []byte) {
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	mtip.inputChan <- fileEntry{
		A: filename,
		B: data,
	}
}

func (mtip *MultiThreadImageProcessor) workerHandl() {
	defer mtip.wg.Done()
	for entry := range mtip.inputChan {
		err := mtip.run(entry.A, entry.B)
		if err != nil {
			slog.Info(
				"failed while running image worker",
				slog.String("filename", entry.A),
				slog.String("error", err.Error()),
			)
			return
		}
	}
}

func (mtip *MultiThreadImageProcessor) run(fileName string, data []byte) (err error) {
	var decodedImg image.Image
	if decodedImg, _, err = image.Decode(bytes.NewReader(data)); err != nil {
		return err
	}

	var finalImgList []image.Image
	if finalImgList, err = mtip.imgPipeline.Process(decodedImg); err != nil {
		return err
	}

	for index, img := range finalImgList {
		err = mtip.fileWriter.Handler(
			fileName+"__"+strconv.Itoa(index)+".jpeg",
			func(writer io.Writer) error {
				return jpeg.Encode(writer, img, &jpeg.Options{Quality: 85})
			},
		)
	}

	return err
}

func (mtip *MultiThreadImageProcessor) Close() error {
	if !mtip.isFinished.Swap(true) {
		close(mtip.inputChan)
	}
	return nil
}

func (mtip *MultiThreadImageProcessor) Shutdown() error {
	err := mtip.Close()
	mtip.wg.Wait() // Wait all goroutines to finish

	err = mtip.fileWriter.Flush()
	return err
}
