package cbxr

import (
	"errors"
	"io"
	"iter"

	"github.com/Jictyvoo/ink_stream/internal/utils"
)

type (
	FileResult        utils.ResultErr[[]byte]
	FileName          string
	FileContentStream interface {
		io.ReaderAt
		io.ReadSeeker
	}
)

type Extractor interface {
	FileSeq() iter.Seq2[FileName, FileResult]
}

var ErrUnsupportedFormat = errors.New("unsupported file format")

func SupportedFileExtensions() []string {
	return []string{".cbz", ".cbr", ".zip", ".rar", ".pdf"}
}
