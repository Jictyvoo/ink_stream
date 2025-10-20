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

func DefaultInsideIgnore() [][2]rune {
	return [][2]rune{{'(', ')'}, {'{', '}'}, {'[', ']'}}
}

func checkIgnores(r int32, stack []rune, ignoreInsideOf [][2]rune, appendDash *bool) []rune {
	foundIdx := slices.IndexFunc(
		ignoreInsideOf, func(toIgnore [2]rune) bool { return toIgnore[0] == r },
	)
	if foundIdx >= 0 {
		stack = append(stack, ignoreInsideOf[foundIdx][1])
		return stack
	}
	// If we are currently inside an ignore section
	if len(stack) > 0 {
		// Check if this rune closes the current section
		if r == stack[len(stack)-1] {
			stack = stack[:len(stack)-1] // pop
			*appendDash = true
		}
		return stack
	}

	return nil
}

func SanitizeName(
	input string, divider rune, runeWriter func(r rune) rune,
	ignoreInsideOf [][2]rune, keepRunes ...rune,
) string {
	if runeWriter == nil {
		runeWriter = func(r rune) rune { return r }
	}
	builder := ([]rune(input))[:0]
	var appendDash bool

	var ignoreStack []rune    // stack to track nested ignores
	for _, r := range input { // If this rune opens an ignore section, push its closing rune
		ignoreStack = checkIgnores(r, ignoreStack, ignoreInsideOf, &appendDash)
		if len(ignoreStack) > 0 {
			continue
		}

		// Handle normal runes
		if unicode.IsLetter(r) || unicode.IsDigit(r) || slices.Contains(keepRunes, r) {
			if appendDash && len(builder) > 0 {
				builder = append(builder, divider)
			}
			builder = append(builder, runeWriter(r))
			appendDash = false
		} else {
			appendDash = true
		}
	}

	return string(builder)
}

func NormalizeName(base string, divider rune, ignoreInsideOf [][2]rune, keepRunes ...rune) string {
	return SanitizeName(base, divider, unicode.ToLower, ignoreInsideOf, keepRunes...)
}
