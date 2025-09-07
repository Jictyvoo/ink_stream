package imgprocessor

import (
	"io"
)

type FileWriter interface {
	Handler(filename string, f func(writer io.Writer) error) error
	Flush() error
}
