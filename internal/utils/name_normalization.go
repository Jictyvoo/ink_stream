package utils

import (
	"path/filepath"
	"strings"
	"unicode"
)

// BuildBaseID creates a stable, HTML-safe base ID from the image source path.
// It keeps ASCII letters and digits, converts others to dashes, collapses
// multiple dashes, trims, and ensures the ID starts with a letter.
func BuildBaseID(imageSrc string) string {
	// Use the filename without extension as the seed if possible
	base := filepath.Base(imageSrc)
	if dotIndex := strings.LastIndexByte(base, '.'); dotIndex > 0 {
		base = base[:dotIndex]
	}
	if base == "." || base == "" {
		base = imageSrc
	}

	// normalize
	var builder strings.Builder
	var appendDash bool
	for _, r := range base {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if appendDash && builder.Len() > 0 {
				builder.WriteRune('-')
			}
			builder.WriteRune(unicode.ToLower(r))
			appendDash = false
		} else {
			appendDash = true
		}
	}
	// trim leading/trailing dashes
	res := builder.String()
	if res == "" {
		res = "img"
	}

	// must start with a letter per HTML4 legacy constraints; for safety we still do it
	if c := res[0]; !(unicode.IsLetter(rune(c))) {
		res = "img-" + res
	}
	return res
}
