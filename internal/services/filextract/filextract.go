package filextract

type FileInfo struct {
	CompleteName string
	BaseName     string
}

type FileOutputWriter interface {
	Close() error
	Shutdown() error
	Process(filename string, data []byte)
}

type FileOutputFactory func(outputDir string) (FileOutputWriter, error)
