package cbxr

import (
	"archive/zip"
	"io"
	"iter"
	"os"
)

type CBZExtractor struct {
	zipReader *zip.Reader
}

func NewCBZExtractor(filePointer *os.File) (*CBZExtractor, error) {
	stat, err := filePointer.Stat()
	if err != nil {
		return nil, err
	}

	var zipFile *zip.Reader
	if zipFile, err = zip.NewReader(filePointer, stat.Size()); err != nil {
		return nil, err
	}

	return &CBZExtractor{zipReader: zipFile}, nil
}

func (e *CBZExtractor) FileSeq() iter.Seq2[FileName, FileResult] {
	return func(yield func(FileName, FileResult) bool) {
		for _, innerFile := range e.zipReader.File {
			if innerFile == nil {
				continue
			}

			yieldResult := FileResult{}
			open, err := innerFile.Open()
			if err != nil {
				yieldResult.Error = err
			} else {
				yieldResult.Data, yieldResult.Error = io.ReadAll(open)
			}

			if !yield(FileName(innerFile.Name), yieldResult) {
				return
			}
		}
	}
}
