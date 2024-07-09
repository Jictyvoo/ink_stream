package main

import (
	"context"
	"github.com/mholt/archiver/v4"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type WriterHandle struct {
	outputDirectory    string
	coverDirectoryName string
}

func (wh WriterHandle) subFolderName(f archiver.File) string {
	directoryName := filepath.Dir(f.NameInArchive)
	if index := strings.LastIndex(directoryName, "(en)"); index >= 0 {
		directoryName = directoryName[0:index]
	}
	directoryName = strings.ReplaceAll(strings.TrimSpace(directoryName), " ", "-")

	if directoryName != f.Name() && directoryName != "." {
		folderDir := filepath.Join(wh.outputDirectory, directoryName)
		if _, err := os.Stat(folderDir); os.IsNotExist(err) {
			if err = os.MkdirAll(folderDir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		return directoryName
	}

	return ""
}

func (wh WriterHandle) handler(_ context.Context, f archiver.File) error {
	reader, err := f.Open()
	if err != nil {
		return err
	}

	filename := f.Name()
	if f.IsDir() || (strings.HasPrefix(strings.ToLower(filename), "cred") && len(filename) >= len("000.jpeg")) {
		return nil
	}

	destinationFolder := wh.outputDirectory
	if strings.Contains(strings.ToLower(filename), ".cover") {
		destinationFolder = wh.coverDirectoryName
	}

	if subFolderName := wh.subFolderName(f); subFolderName != "" {
		destinationFolder = filepath.Join(wh.outputDirectory, subFolderName)
	}

	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	writeFile, err := os.Create(destinationFolder + "/" + strings.TrimLeft(filename, "."))
	if err != nil {
		return err
	}
	defer writeFile.Close()
	_, err = writeFile.Write(data)

	return err
}
