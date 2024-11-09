package filextract

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
	"github.com/Jictyvoo/ink_stream/internal/utils"
)

type (
	fileEntry                 utils.Entry2[string, []byte]
	MultiThreadImageProcessor struct {
		fileWriter    outdirwriter.WriterHandle
		imgPipeline   imageparser.ImagePipeline
		wg            sync.WaitGroup
		inputChan     chan fileEntry
		numGoroutines uint8
		isFinished    atomic.Bool
	}
)

func NewMultiThreadImageProcessor(
	extractDir string,
	imgPipeline imageparser.ImagePipeline,
) *MultiThreadImageProcessor {
	mtip := &MultiThreadImageProcessor{
		fileWriter:    outdirwriter.NewWriterHandle(extractDir),
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
			log.Printf("failed while running worker for %s: %s\n", entry.A, err.Error())
			return
		}
	}
}

func (mtip *MultiThreadImageProcessor) run(fileName string, data []byte) (err error) {
	var decodedImg, finalImg image.Image
	if decodedImg, _, err = image.Decode(bytes.NewReader(data)); err != nil {
		return err
	}

	finalImg, err = mtip.imgPipeline.Process(decodedImg)
	if err != nil {
		return
	}

	err = mtip.fileWriter.Handler(fileName+".jpeg", func(writer io.Writer) error {
		return jpeg.Encode(writer, finalImg, &jpeg.Options{Quality: 85})
	})

	return
}

func (mtip *MultiThreadImageProcessor) Close() error {
	if !mtip.isFinished.Swap(true) {
		close(mtip.inputChan)
	}
	return nil
}

func (mtip *MultiThreadImageProcessor) Shutdown() error {
	mtip.wg.Wait() // Wait all goroutines to finish

	err := mtip.fileWriter.OnFinish()
	return err
}
