package extractor

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const coverDirSuffix = "0000_Cover"

func fileIsCover(filename string) bool {
	return strings.Contains(strings.ToLower(filename), ".cover")
}

func (fp *FileProcessorWorker) createOutDir(extractDir string, suffix string) error {
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	// Create a covers output directory
	if err := os.MkdirAll(filepath.Join(extractDir, suffix), 0755); err != nil {
		log.Printf("Failed to create directory for extraction: %v", err)
		return err
	}

	return nil
}
