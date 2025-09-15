package utils

import (
	"reflect"
	"slices"
	"testing"
)

func TestCollapseFilesByExt(t *testing.T) {
	testCases := []struct {
		name     string
		files    []string
		exts     []string
		expected []string
	}{
		{
			name: "all files in folder allowed",
			files: []string{
				"/project/images/cat.jpg",
				"/project/images/dog.png",
			},
			exts:     []string{".jpg", ".png"},
			expected: []string{"/project/images"},
		},
		{
			name: "folder with mixed extensions",
			files: []string{
				"/project/images/cat.jpg",
				"/project/images/readme.txt",
			},
			exts: []string{".jpg", ".png"},
			expected: []string{
				"/project/images/cat.jpg",
				"/project/images/readme.txt",
			},
		},
		{
			name: "multiple folders, one collapses, one not",
			files: []string{
				"/project/images/cat.jpg",
				"/project/images/dog.png",
				"/project/docs/readme.md",
				"/project/docs/manual.pdf",
			},
			exts: []string{".jpg", ".png"},
			expected: []string{
				"/project/images",
				"/project/docs/readme.md",
				"/project/docs/manual.pdf",
			},
		},
		{
			name: "extension case insensitivity",
			files: []string{
				"/project/images/cat.JPG",
				"/project/images/dog.PNG",
			},
			exts:     []string{".jpg", ".png"},
			expected: []string{"/project/images"},
		},
		{
			name:     "empty file list",
			files:    []string{},
			exts:     []string{".jpg"},
			expected: []string{},
		},
		{
			name: "nested subdirectories with only allowed extensions collapse to root",
			files: []string{
				"/project/gallery/sub1/cat.jpg",
				"/project/gallery/sub1/dog.png",
				"/project/gallery/sub2/bird.jpg",
				"/project/gallery/sub2/fish.png",
				"/project/gallery/thumbnail.png",
			},
			exts:     []string{".jpg", ".png"},
			expected: []string{"/project/gallery"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := CollapseFilesByExt(tt.files, tt.exts)

			// order of result doesn’t matter
			slices.Sort(got)
			slices.Sort(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
