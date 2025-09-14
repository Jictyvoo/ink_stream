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
			want:     "logo",
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
