package cbxr

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

func TestFolderExtractor_FileSeq(t *testing.T) {
	// Create test files
	testFiles := map[string]string{
		"image1.jpg":         "jpg content",
		"image2.png":         "png content",
		"document.pdf":       "pdf content",
		"readme.txt":         "txt content",
		"image3.gif":         "gif content",
		"subdir/image4.jpeg": "jpeg content",
	}

	t.Run("basic iteration", func(t *testing.T) {
		// Create a temporary directory structure for testing
		tempDir := t.TempDir()

		for filename, content := range testFiles {
			path := filepath.Join(tempDir, filename)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				t.Fatalf("Failed to write file: %v", err)
			}
		}

		// Create FolderExtractor
		folder, err := os.Open(tempDir)
		if err != nil {
			t.Fatalf("Failed to open folder: %v", err)
		}
		defer folder.Close()

		extractor, newExtractorErr := NewFolderExtractor(folder)
		if newExtractorErr != nil {
			t.Fatalf("Failed to create FolderExtractor: %v", newExtractorErr.Error())
		}

		count := 0
		seenFiles := make(map[string]bool)
		for filename, result := range extractor.FileSeq() {
			count++
			seenFiles[string(filename)] = true
			if result.Error != nil {
				t.Errorf("Unexpected error for file %s: %v", filename, result.Error)
			}
			if result.Data == nil {
				t.Errorf("Expected data for file %s, got nil", filename)
			}
		}
		// Check that we found expected files
		expectedFiles := make([]string, 0, len(testFiles))
		for filename := range maps.Keys(testFiles) {
			if slices.Contains(imgutils.SupportedImageFormats(), filepath.Ext(filename)) {
				expectedFiles = append(expectedFiles, filepath.Join(tempDir, filename))
			}
		}
		for _, expectedFile := range expectedFiles {
			if !seenFiles[expectedFile] {
				t.Errorf("Expected file %s was not found in iteration", expectedFile)
			}
		}
		// Should have found 4 files (jpg, png, gif, jpeg)
		if count != 4 {
			t.Errorf("Expected 4 files, got %d", count)
		}
	})
}
