package filextract

type FileInfo struct {
	CompleteName string
	BaseName     string
}

type CropStyle uint8

const (
	CropBasic CropStyle = iota
	CropNormal
	CropAggressive
)

type Options struct{}
