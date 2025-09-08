package imgprocessor

import (
	"io"
)

type WriterCallback func(writer io.Writer) error

type FileWriter interface {
	Handler(filename string, f WriterCallback) error
	Flush() error
}
