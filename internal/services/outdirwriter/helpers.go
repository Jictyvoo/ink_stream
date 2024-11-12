package outdirwriter

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const CoverDirSuffix = "0000_Cover"

func fileIsCover(filename string) bool {
	return strings.Contains(strings.ToLower(filename), ".cover")
}

func CreateOutDir(extractDir string, suffix string) error {
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		slog.Error("Failed to create directory for extraction", slog.String("error", err.Error()))
		return err
	}

	// Create a covers output directory
	if err := os.MkdirAll(filepath.Join(extractDir, suffix), 0755); err != nil {
		slog.Error("Failed to create directory for extraction", slog.String("error", err.Error()))
		return err
	}

	return nil
}
