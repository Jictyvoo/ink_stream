package outdirwriter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jictyvoo/ink_stream/internal/utils"
)

const defaultContentDir = "content-main"

type (
	folderInfoCounter struct{ onRoot, cover, onSubDir, total uint32 }
	WriterHandle      struct {
		outputDirectory    string
		coverDirectoryName string
		folderCounter      *folderInfoCounter
	}
)

func NewWriterHandle(extractDir string) WriterHandle {
	return WriterHandle{
		outputDirectory:    extractDir,
		coverDirectoryName: filepath.Join(extractDir, CoverDirSuffix),
		folderCounter:      &folderInfoCounter{},
	}
}

func (wh WriterHandle) defaultDir() string {
	return filepath.Join(wh.outputDirectory, defaultContentDir)
}

func (wh WriterHandle) subFolderName(absFilename string) (directoryName string) {
	directoryName = filepath.Dir(absFilename)
	if index := strings.LastIndex(directoryName, "(en)"); index >= 0 {
		directoryName = directoryName[0:index]
	}
	directoryName = strings.ReplaceAll(strings.TrimSpace(directoryName), " ", "-")

	filename := filepath.Base(absFilename)
	folderDir := wh.defaultDir()
	if directoryName != filename && directoryName != "." {
		folderDir = filepath.Join(wh.outputDirectory, directoryName)
		wh.folderCounter.onSubDir++
	} else {
		directoryName = defaultContentDir
		switch fileIsCover(filename) {
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

func (wh WriterHandle) Handler(filename string, data []byte) error {
	destinationFolder := wh.outputDirectory
	if fileIsCover(filename) {
		destinationFolder = wh.coverDirectoryName
		wh.folderCounter.cover++
	} else if subFolderName := wh.subFolderName(filename); subFolderName != "" {
		destinationFolder = filepath.Join(wh.outputDirectory, subFolderName)
	}

	writeFile, err := os.Create(destinationFolder + "/" + strings.TrimLeft(filename, "."))
	if err != nil {
		return err
	}
	defer writeFile.Close()
	_, err = writeFile.Write(data)

	return err
}

func (wh WriterHandle) OnFinish() (err error) {
	if wh.folderCounter.onRoot > 0 && wh.folderCounter.onSubDir > 0 {
		if wh.folderCounter.cover > 0 {
			err = fmt.Errorf("only one of onRoot and onSubDir may be specified")
			return err
		}
		// Move all content-main to work as _0Cover
		err = errors.Join(
			os.Remove(wh.coverDirectoryName),
			os.Rename(wh.defaultDir(), wh.coverDirectoryName),
		)
	}

	return
}
