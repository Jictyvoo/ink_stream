package bootstrap

import (
	"io"

	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

type FileWriterWrapper struct {
	outdirwriter.WriterHandle
}

func NewFileWriterWrapper(extractDir string) (*FileWriterWrapper, error) {
	wh, err := outdirwriter.NewWriterHandle(extractDir)
	return &FileWriterWrapper{WriterHandle: wh}, err
}

func (f FileWriterWrapper) Close() error {
	return nil
}

func (f FileWriterWrapper) Shutdown() error {
	return f.WriterHandle.Flush()
}

func (f FileWriterWrapper) Process(filename string, data []byte) {
	_ = f.WriterHandle.Handler(
		filename, func(writer io.Writer) (inktypes.ImageMetadata, error) {
			_, err := writer.Write(data)
			return inktypes.ImageMetadata{}, err
		},
	)
}
