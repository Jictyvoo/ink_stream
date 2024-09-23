package extractor

import (
	"context"
	"io"
	"io/fs"
	"strings"

	"github.com/mholt/archiver/v4"
)

const coverDirSuffix = "0000_Cover"

func checkFileFormat(filename string, file io.Reader) (io.Reader, archiver.Extractor, error) {
	format, fileReader, err := archiver.Identify(filename, file)
	if err != nil {
		return nil, nil, err
	}

	// want to extract something?
	if ex, ok := format.(archiver.Extractor); ok {
		return fileReader, ex, nil
	}

	// or maybe it's compressed and you want to decompress it?
	if decom, ok := format.(archiver.Decompressor); ok {
		var rc io.ReadCloser
		if rc, err = decom.OpenReader(fileReader); err != nil {
			return nil, nil, err
		}
		defer rc.Close()
		return checkFileFormat(filename, rc)
	}
	return nil, nil, nil
}

func getAllNames(filename string) (result []string) {
	fsys, err := archiver.FileSystem(context.Background(), filename)
	if err != nil {
		return nil
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		result = append(result, path)
		return nil
	})

	return
}

func fileIsCover(filename string) bool {
	return strings.Contains(strings.ToLower(filename), ".cover")
}
