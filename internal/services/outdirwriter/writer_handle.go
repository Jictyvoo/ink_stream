package outdirwriter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/Jictyvoo/ink_stream/internal/services/imgprocessor"
	"github.com/Jictyvoo/ink_stream/internal/utils"
	"github.com/Jictyvoo/ink_stream/pkg/inktypes"
)

const defaultContentDir = "content-main"

type (
	folderInfoCounter struct{ onRoot, cover, onSubDir, total atomic.Uint32 }
	WriterHandle      struct {
		outputDirectory    string
		coverDirectoryName string
		folderCounter      *folderInfoCounter
	}
)

func NewWriterHandle(extractDir string) (WriterHandle, error) {
	// Create the directory for the extracted files
	if err := CreateOutDir(extractDir, CoverDirSuffix); err != nil {
		return WriterHandle{}, err
	}

	wh := WriterHandle{
		outputDirectory:    extractDir,
		coverDirectoryName: filepath.Join(extractDir, CoverDirSuffix),
		folderCounter:      &folderInfoCounter{},
	}
	return wh, nil
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
		wh.folderCounter.onSubDir.Add(1)
	} else {
		directoryName = defaultContentDir
		switch fileIsCover(filename) {
		case true:
			wh.folderCounter.cover.Add(1)
		case false:
			wh.folderCounter.onRoot.Add(1)
		}
	}

	if err := utils.CreateDirIfNotExist(folderDir); err != nil {
		log.Fatal(err)
	}

	wh.folderCounter.total.Add(1)
	return directoryName
}

// ExecuteFileWrite writes the file and also returns the
// image metadata and the absolute path of the written file.
func (wh WriterHandle) ExecuteFileWrite(
	filename string, callback imgprocessor.WriterCallback,
) (meta inktypes.ImageMetadata, absPath string, err error) {
	destinationFolder := wh.outputDirectory
	if fileIsCover(filename) {
		destinationFolder = wh.coverDirectoryName
		wh.folderCounter.cover.Add(1)
	} else if subFolderName := wh.subFolderName(filename); subFolderName != "" {
		destinationFolder = filepath.Join(wh.outputDirectory, subFolderName)
	}

	absPath = filepath.Join(destinationFolder, strings.TrimLeft(filepath.Base(filename), "."))
	writeFile, creatErr := os.Create(absPath)
	if creatErr != nil {
		return inktypes.ImageMetadata{}, "", creatErr
	}
	defer writeFile.Close()

	meta, err = callback(writeFile)
	return meta, absPath, err
}

func (wh WriterHandle) Handler(filename string, callback imgprocessor.WriterCallback) error {
	_, _, err := wh.ExecuteFileWrite(filename, callback)
	return err
}

func (wh WriterHandle) Flush() (err error) {
	if wh.folderCounter.onRoot.Load() > 0 && wh.folderCounter.onSubDir.Load() > 0 {
		if wh.folderCounter.cover.Load() > 0 {
			err = fmt.Errorf("only one of onRoot and onSubDir may be specified")
			return err
		}
		// Move all content-main to work as _0Cover
		err = errors.Join(
			os.Remove(wh.coverDirectoryName),
			os.Rename(wh.defaultDir(), wh.coverDirectoryName),
		)
	}

	// After finishing file processing, start the post-analysis
	if err = MoveFirstFileToCoverFolder(wh.outputDirectory); err != nil {
		return err
	}
	return err
}
