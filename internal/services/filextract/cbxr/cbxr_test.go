package cbxr

import (
	"bytes"
	"iter"
	"slices"
	"testing"
)

type fileSeqProvider interface {
	FileSeq() iter.Seq2[FileName, FileResult]
}

type ExtractTestSuite struct {
	WantErr           bool
	WantFiles         []string
	WantFileData      map[string][]byte
	ExpectedCount     *int // optional: if provided, enforce total yielded files
	RequireDataNonNil bool // if true, fail when any yielded Data is nil
}

func (s ExtractTestSuite) Run(t *testing.T, extractor fileSeqProvider) {
	t.Helper()

	var files []string
	fileData := map[string][]byte{}
	collectData := s.WantFileData != nil || s.RequireDataNonNil

	for name, result := range extractor.FileSeq() {
		files = append(files, string(name))
		if result.Error != nil && !s.WantErr {
			t.Errorf("Unexpected error: %v", result.Error)
		}
		if collectData {
			fileData[string(name)] = result.Data
			if s.RequireDataNonNil && result.Data == nil {
				t.Errorf("Expected data for file %s, got nil", name)
			}
		}
	}

	if s.ExpectedCount != nil {
		if len(files) != *s.ExpectedCount {
			t.Errorf("Expected %d files, got %d", *s.ExpectedCount, len(files))
		}
	}

	if s.WantFiles != nil {
		if len(files) != len(s.WantFiles) && s.ExpectedCount == nil {
			// Only enforce exact length when count expectation wasn't explicitly provided
			t.Errorf("Expected %d files, got %d", len(s.WantFiles), len(files))
		}
		for _, wantFile := range s.WantFiles {
			if !slices.Contains(files, wantFile) {
				t.Errorf("Expected file %q not found", wantFile)
			}
		}
	}

	if s.WantFileData != nil {
		for fileName, wantData := range s.WantFileData {
			gotData, exists := fileData[fileName]
			if !exists {
				t.Errorf("Expected data for file %q not found", fileName)
				continue
			}
			if !bytes.Equal(gotData, wantData) {
				t.Errorf("Data mismatch for file %q: got %v, want %v", fileName, gotData, wantData)
			}
		}
	}
}

func TestSupportedFileExtensions(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "default extensions",
			want: []string{".cbz", ".cbr", ".zip", ".rar", ".pdf"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SupportedFileExtensions(); !slices.Equal(tt.want, got) {
				t.Errorf("SupportedFileExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}
