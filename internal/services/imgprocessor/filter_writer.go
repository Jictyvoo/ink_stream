package imgprocessor

import (
	"io"

	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

type WriterCallback func(writer io.Writer) (metadata inktypes.ImageMetadata, err error)

type FileWriter interface {
	Handler(filename string, f WriterCallback) error
	Flush() error
}
