package inktypes

import "testing"

func TestNewImageEncodingOptions(t *testing.T) {
	tests := []struct {
		name     string
		quality  uint8
		format   ImageFormat
		expected ImageEncodingOptions
	}{
		{
			name:    "valid quality and format",
			quality: 80,
			format:  FormatJPEG,
			expected: ImageEncodingOptions{
				Quality: 80,
				Format:  FormatJPEG,
			},
		},
		{
			name:    "quality over 100 should be capped at 100",
			quality: 150,
			format:  "png",
			expected: ImageEncodingOptions{
				Quality: 100,
				Format:  "png",
			},
		},
		{
			name:    "zero quality should be 60",
			quality: 0,
			format:  "webp",
			expected: ImageEncodingOptions{
				Quality: 60,
				Format:  "webp",
			},
		},
		{
			name:    "format conversion to lowercase",
			quality: 90,
			format:  "JPEG",
			expected: ImageEncodingOptions{
				Quality: 90,
				Format:  FormatJPEG,
			},
		},
		{
			name:    "empty format",
			quality: 75,
			format:  "",
			expected: ImageEncodingOptions{
				Quality: 75,
				Format:  FormatJPEG,
			},
		},
		{
			name:    "mixed case format",
			quality: 60,
			format:  "Jpeg",
			expected: ImageEncodingOptions{
				Quality: 60,
				Format:  FormatJPEG,
			},
		},
		{
			name:    "quality exactly 100",
			quality: 100,
			format:  "gif",
			expected: ImageEncodingOptions{
				Quality: 100,
				Format:  "gif",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewImageEncodingOptions(tt.quality, tt.format)
			if result.Quality != tt.expected.Quality {
				t.Errorf("Quality mismatch: got %d, want %d", result.Quality, tt.expected.Quality)
			}
			if result.Format != tt.expected.Format {
				t.Errorf("Format mismatch: got %s, want %s", result.Format, tt.expected.Format)
			}
		})
	}
}
