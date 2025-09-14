package main

import (
	"io"

	"github.com/Jictyvoo/ink_stream/internal/services/outdirwriter"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

type FileWriterWrapper struct {
	outdirwriter.WriterHandle
}

func NewFileWriterWrapper(extractDir string) *FileWriterWrapper {
	return &FileWriterWrapper{
		outdirwriter.NewWriterHandle(extractDir),
	}
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
