package imgutils

import (
	"image"
	"iter"
)

// IsImageEqual returns true if two images a and b are equal.
func IsImageEqual(a, b image.Image) bool {
	if !a.Bounds().Eq(b.Bounds()) {
		return false
	}

	for y := 0; y < a.Bounds().Dy(); y++ {
		for x := 0; x < a.Bounds().Dx(); x++ {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	return true
}

func Iterator(input image.Image) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		bounds := input.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if !yield(x, y) {
					return
				}
			}
		}
	}
}
