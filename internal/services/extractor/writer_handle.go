package extractor

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"

	"github.com/Jictyvoo/ink_stream/internal/utils"
)

type (
	folderInfoCounter struct{ onRoot, cover, onSubDir, total uint32 }
	WriterHandle      struct {
		outputDirectory    string
		coverDirectoryName string
		folderCounter      *folderInfoCounter
	}
)

const defaultContentDir = "content-main"

func (wh WriterHandle) defaultDir() string {
	return filepath.Join(wh.outputDirectory, defaultContentDir)
}

func (wh WriterHandle) subFolderName(f archiver.File) (directoryName string) {
	directoryName = filepath.Dir(f.NameInArchive)
	if index := strings.LastIndex(directoryName, "(en)"); index >= 0 {
		directoryName = directoryName[0:index]
	}
	directoryName = strings.ReplaceAll(strings.TrimSpace(directoryName), " ", "-")

	folderDir := wh.defaultDir()
	if directoryName != f.Name() && directoryName != "." {
		folderDir = filepath.Join(wh.outputDirectory, directoryName)
		wh.folderCounter.onSubDir++
	} else {
		directoryName = defaultContentDir
		switch fileIsCover(filepath.Base(f.NameInArchive)) {
		case true:
			wh.folderCounter.cover++
		case false:
			wh.folderCounter.onRoot++
		}
	}

	if err := utils.CreateDirIfNotExist(folderDir); err != nil {
		log.Fatal(err)
	}

	wh.folderCounter.total++
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
