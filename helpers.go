package main

import (
	"context"
	"github.com/mholt/archiver/v4"
	"io"
	"io/fs"
)

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
		rc, err := decom.OpenReader(fileReader)
		if err != nil {
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

		result = append(result, path)
		return nil
	})
	return
}
