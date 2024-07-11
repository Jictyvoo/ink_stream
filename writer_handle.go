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

type (
	folderInfoCounter struct{ onRoot, cover, onSubDir, total uint32 }
	WriterHandle      struct {
		outputDirectory    string
		coverDirectoryName string
		folderCounter      *folderInfoCounter
	}
)

func (wh WriterHandle) subFolderName(f archiver.File) (directoryName string) {
	directoryName = filepath.Dir(f.NameInArchive)
	if index := strings.LastIndex(directoryName, "(en)"); index >= 0 {
		directoryName = directoryName[0:index]
	}
	directoryName = strings.ReplaceAll(strings.TrimSpace(directoryName), " ", "-")

	const defaultContentDir = "content-main"
	folderDir := filepath.Join(wh.outputDirectory, defaultContentDir)
	if directoryName != f.Name() && directoryName != "." {
		folderDir = filepath.Join(wh.outputDirectory, directoryName)
		wh.folderCounter.onSubDir++
	} else {
		directoryName = defaultContentDir
	}

	if _, err := os.Stat(folderDir); os.IsNotExist(err) {
		if err = os.MkdirAll(folderDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	wh.folderCounter.onRoot++
	return
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
	if fileIsCover(filename) {
		destinationFolder = wh.coverDirectoryName
		wh.folderCounter.cover++
	} else if subFolderName := wh.subFolderName(f); subFolderName != "" {
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
