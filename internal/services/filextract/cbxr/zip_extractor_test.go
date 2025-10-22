package cbxr

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
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

			// Create CBZExtractor
			zipFile, openErr := os.Open(filename)
			if openErr != nil {
				t.Fatalf("Failed to open zip file: %v", openErr)
			}
			defer zipFile.Close()
			extractor, newExtractorErr := NewCBZExtractor(zipFile)
			if newExtractorErr != nil {
				t.Fatalf("Failed to create extractor: %v", newExtractorErr.Error())
			}

			suite := ExtractTestSuite{
				WantErr:      tt.wantErr,
				WantFiles:    tt.wantFiles,
				WantFileData: tt.wantFileData,
			}
			suite.Run(t, extractor)
		})
	}
}
