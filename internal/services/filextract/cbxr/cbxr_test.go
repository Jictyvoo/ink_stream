package cbxr

import (
	"slices"
	"testing"
)

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
