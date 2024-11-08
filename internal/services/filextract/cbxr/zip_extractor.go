package cbxr

import (
	"archive/zip"
	"io"
	"iter"
)

type CBZExtractor struct {
	zipReader *zip.Reader
}

func NewCBZExtractor(fileReader FileContentStream, size int64) (*CBZExtractor, error) {
	zipFile, err := zip.NewReader(fileReader, size)
	if err != nil {
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
