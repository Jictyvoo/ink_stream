package utils

import (
	"path/filepath"
	"slices"
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
		base = filepath.Join(filepath.Dir(imageSrc), base[:dotIndex])
	}
	if base == "." || base == "" {
		base = imageSrc
	}

	// normalize
	res := NormalizeName(base, '-', nil)
	if res == "" {
		res = "img"
	}

	// must start with a letter per HTML4 legacy constraints; for safety we still do it
	if c := res[0]; !(unicode.IsLetter(rune(c))) {
		res = "img-" + res
	}
	return res
}

func NormalizeName(base string, divider rune, ignoreInsideOf [][2]rune, keepRunes ...rune) string {
	var builder strings.Builder
	var appendDash bool
	var expectToCloseOpened rune // for cases to check for a rune to close like parenthesis
	for _, r := range base {
		if expectToCloseOpened != 0 && r != expectToCloseOpened {
			continue
		}
		expectToCloseOpened = 0
		if unicode.IsLetter(r) || unicode.IsDigit(r) || slices.Contains(keepRunes, r) {
			if appendDash && builder.Len() > 0 {
				builder.WriteRune(divider)
			}
			builder.WriteRune(unicode.ToLower(r))
			appendDash = false
		} else {
			foundIdx := slices.IndexFunc(ignoreInsideOf, func(toIgnore [2]rune) bool { return toIgnore[0] == r })
			if foundIdx >= 0 {
				expectToCloseOpened = ignoreInsideOf[foundIdx][1]
			}
			appendDash = true
		}
	}
	res := builder.String()
	return res
}
