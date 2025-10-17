package utils

import (
	"testing"
)

func TestBuildBaseID(t *testing.T) {
	testCases := []struct {
		name     string
		imageSrc string
		want     string
	}{
		{
			name:     "Simple filename",
			imageSrc: "photo.jpg",
			want:     "photo",
		},
		{
			name:     "Uppercase letters",
			imageSrc: "MyImage.PNG",
			want:     "myimage",
		},
		{
			name:     "Spaces and special chars",
			imageSrc: "hello world!.png",
			want:     "hello-world",
		},
		{
			name:     "Multiple dashes collapse",
			imageSrc: "a---b---c.png",
			want:     "a-b-c",
		},
		{
			name:     "Starts with digit",
			imageSrc: "123image.gif",
			want:     "img-123image",
		},
		{
			name:     "Only special chars",
			imageSrc: "@#$%^&*()",
			want:     "img",
		},
		{
			name:     "Empty string",
			imageSrc: "",
			want:     "img",
		},
		{
			name:     "Path with directories",
			imageSrc: "/var/images/icons/logo.svg",
			want:     "var-images-icons-logo",
		},
		{
			name:     "Dot as filename",
			imageSrc: ".",
			want:     "img",
		},
		{
			name:     "No extension",
			imageSrc: "avatar",
			want:     "avatar",
		},
		{
			name:     "Filename with unicode letters",
			imageSrc: "café.png",
			want:     "café",
		},
		{
			name:     "Filename with unicode symbols",
			imageSrc: "猫🐱.jpg",
			want:     "猫",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildBaseID(tt.imageSrc)
			if got != tt.want {
				t.Errorf("BuildBaseID(%q) = %q, want %q", tt.imageSrc, got, tt.want)
			}
		})
	}
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name           string
		base           string
		divider        rune
		ignoreInsideOf [][2]rune
		keepRunes      []rune
		want           string
	}{
		{
			name:    "basic letters and digits",
			base:    "Hello World 123",
			divider: '-',
			want:    "hello-world-123",
		},
		{
			name:    "multiple dividers collapse",
			base:    "Hello---World",
			divider: '-',
			want:    "hello-world",
		},
		{
			name:      "keep underscore rune",
			base:      "Hello_World!",
			divider:   '-',
			keepRunes: []rune{'_'},
			want:      "hello_world",
		},
		{
			name:    "ignore inside parentheses",
			base:    "Hello (World) Test",
			divider: '-',
			ignoreInsideOf: [][2]rune{
				{'(', ')'},
			},
			want: "hello-test",
		},
		{
			name:    "ignore inside brackets",
			base:    "Hello [Internal] Name",
			divider: '-',
			ignoreInsideOf: [][2]rune{
				{'[', ']'},
			},
			want: "hello-name",
		},
		{
			name:    "ignore inside curly braces and brackets together",
			base:    "A{Skip}B[Hide]C",
			divider: '-',
			ignoreInsideOf: [][2]rune{
				{'{', '}'},
				{'[', ']'},
			},
			want: "a-b-c",
		},
		{
			name:    "mixed special characters",
			base:    "My.Name/Is!InkStream",
			divider: '-',
			want:    "my-name-is-inkstream",
		},
		{
			name:    "leading and trailing special chars",
			base:    "--Hello--",
			divider: '-',
			want:    "hello",
		},
		{
			name:    "nested ignored sections (outer only respected)",
			base:    "Outer(Start(Inner)End)After",
			divider: '-',
			ignoreInsideOf: [][2]rune{
				{'(', ')'},
			},
			want: "outer-after",
		},
		{
			name:    "empty string",
			base:    "",
			divider: '-',
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeName(tt.base, tt.divider, tt.ignoreInsideOf, tt.keepRunes...)
			if got != tt.want {
				t.Errorf("NormalizeName(%q) = %q, want %q", tt.base, got, tt.want)
			}
		})
	}
}
