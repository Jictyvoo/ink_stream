package cbxr

import (
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

type FolderExtractor struct {
	folderPointer *os.File
}

func NewFolderExtractor(folderPointer *os.File) (*FolderExtractor, error) {
	if folderPointer == nil {
		return nil, fmt.Errorf("folderPointer is nil")
	}

	fileStat, err := folderPointer.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get folderPointer stat: %w", err)
	}
	if !fileStat.IsDir() {
		return nil, fmt.Errorf("folderPointer is not a directory")
	}
	return &FolderExtractor{folderPointer: folderPointer}, nil
}

func (e *FolderExtractor) FileSeq() iter.Seq2[FileName, FileResult] {
	return func(yield func(FileName, FileResult) bool) {
		root := e.folderPointer.Name()
		supportedFormats := imgutils.SupportedImageFormats()

		_ = filepath.WalkDir(root, func(path string, dirEntry os.DirEntry, err error) error {
			if err != nil {
				yield("", FileResult{Error: err})
				return filepath.SkipDir
			}
			if dirEntry.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(dirEntry.Name()))
			if !slices.Contains(supportedFormats, ext) {
				return nil
			}

			var result FileResult
			result.Data, result.Error = os.ReadFile(path)
			if !yield(FileName(path), result) {
				return io.EOF // stop walking early
			}
			return nil
		})
	}
}
