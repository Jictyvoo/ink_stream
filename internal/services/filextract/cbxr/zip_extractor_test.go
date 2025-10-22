package cbxr

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestCBZExtractor_FileSeq(t *testing.T) {
	tests := []struct {
		name         string
		files        map[string]string // path -> content
		wantErr      bool
		wantFiles    []string
		wantFileData map[string][]byte
	}{
		{
			name:      "empty zip",
			files:     map[string]string{},
			wantFiles: []string{},
		},
		{
			name:         "single file",
			files:        map[string]string{"file1.txt": "content1"},
			wantFiles:    []string{"file1.txt"},
			wantFileData: map[string][]byte{"file1.txt": []byte("content1")},
		},
		{
			name: "multiple files",
			files: map[string]string{
				"file1.txt": "content1",
				"file2.txt": "content2",
			},
			wantFiles: []string{"file1.txt", "file2.txt"},
			wantFileData: map[string][]byte{
				"file1.txt": []byte("content1"),
				"file2.txt": []byte("content2"),
			},
		},
		{
			name:         "file with special characters",
			files:        map[string]string{"file with spaces.txt": "content"},
			wantFiles:    []string{"file with spaces.txt"},
			wantFileData: map[string][]byte{"file with spaces.txt": []byte("content")},
		},
		{
			name:         "nested directories",
			files:        map[string]string{"dir1/dir2/file.txt": "content"},
			wantFiles:    []string{"dir1/dir2/file.txt"},
			wantFileData: map[string][]byte{"dir1/dir2/file.txt": []byte("content")},
		},
		{
			name:         "file with zero size",
			files:        map[string]string{"empty.txt": ""},
			wantFiles:    []string{"empty.txt"},
			wantFileData: map[string][]byte{"empty.txt": {}},
		},
		{
			name:         "large file",
			files:        map[string]string{"large.txt": strings.Repeat("a", 10000)},
			wantFiles:    []string{"large.txt"},
			wantFileData: map[string][]byte{"large.txt": []byte(strings.Repeat("a", 10000))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the zip file content inside the test run using tt.files
			zipBuffer := new(bytes.Buffer)
			writer := zip.NewWriter(zipBuffer)
			for path, content := range tt.files {
				fw, err := writer.Create(path)
				if err != nil {
					t.Fatalf("Failed to create file %q in zip: %v", path, err)
				}
				if _, err = fw.Write([]byte(content)); err != nil {
					t.Fatalf("Failed to write content for %q in zip: %v", path, err)
				}
			}
			if err := writer.Close(); err != nil {
				t.Fatalf("Failed to close zip writer: %v", err)
			}

			// Create a temporary directory structure for testing
			tempDir := t.TempDir()
			filename := filepath.Join(tempDir, tt.name+".zip")
			if err := os.WriteFile(filename, zipBuffer.Bytes(), 0o644); err != nil {
				t.Fatalf("Failed to write file: %v", err)
			}

			// Create FolderExtractor
			folder, openErr := os.Open(filename)
			if openErr != nil {
				t.Fatalf("Failed to open folder: %v", openErr)
			}
			defer folder.Close()
			extractor, newExtractorErr := NewCBZExtractor(folder)
			if newExtractorErr != nil {
				t.Fatalf("Failed to create extractor: %v", newExtractorErr.Error())
			}

			var files []string
			var fileData map[string][]byte
			if tt.wantFileData != nil {
				fileData = make(map[string][]byte)
			}

			for name, result := range extractor.FileSeq() {
				files = append(files, string(name))
				if tt.wantFileData != nil {
					fileData[string(name)] = result.Data
				}
				if result.Error != nil && !tt.wantErr {
					t.Errorf("Unexpected error: %v", result.Error)
				}
			}

			if len(files) != len(tt.wantFiles) {
				t.Errorf("Expected %d files, got %d", len(tt.wantFiles), len(files))
			}
			for _, wantFile := range tt.wantFiles {
				if !slices.Contains(files, wantFile) {
					t.Errorf("Expected file %q not found", wantFile)
				}
			}

			if tt.wantFileData != nil {
				for fileName, wantData := range tt.wantFileData {
					gotData, exists := fileData[fileName]
					if !exists {
						t.Errorf("Expected data for file %q not found", fileName)
						continue
					}
					if !bytes.Equal(gotData, wantData) {
						t.Errorf(
							"Data mismatch for file %q: got %v, want %v",
							fileName, gotData, wantData,
						)
					}
				}
			}
		})
	}
}
