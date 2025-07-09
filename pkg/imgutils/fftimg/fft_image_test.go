package fftimg

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func almostEqual(a, b float64, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestFFTImage_GrayscaleFFT(t *testing.T) {
	// 6x6 grayscale image with a simple pattern
	w, h := 6, 6
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Fill with a pattern (for reproducibility, use increasing values)
	vals := [][]uint8{
		{0, 1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5, 6},
		{2, 3, 4, 5, 6, 7},
		{3, 4, 5, 6, 7, 8},
		{4, 5, 6, 7, 8, 9},
		{5, 6, 7, 8, 9, 10},
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := vals[y][x]
			img.Set(x, y, color.RGBA{v, v, v, 255})
		}
	}

	fftImg := NewFFTImage(img)
	ffts := fftImg.FFT()
	// Only test the R channel (since it's grayscale, all channels are the same)
	mag := make([][]float64, h)
	for y := 0; y < h; y++ {
		mag[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			mag[y][x] = math.Round(cmplxAbs(ffts[0][y][x])*10) / 10 // round to 1 decimal
		}
	}

	expected := [][]float64{
		{40, 0.4, 0.5, 1.4, 3.2, 1.1},
		{0.4, 0.5, 0.7, 1.8, 4, 1.2},
		{0.5, 0.7, 1.1, 2.8, 5.9, 1.7},
		{1.4, 1.8, 2.8, 6.8, 14.1, 3.8},
		{3.2, 4, 5.9, 14.1, 27.5, 6.8},
		{1.1, 1.2, 1.7, 3.8, 6.8, 1.6},
	}

	tol := 0.2 // allow some tolerance for floating point differences
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !almostEqual(mag[y][x], expected[y][x], tol) {
				t.Errorf("FFT magnitude mismatch at (%d,%d): got %v, want %v", x, y, mag[y][x], expected[y][x])
			}
		}
	}
}

func cmplxAbs(c complex128) float64 {
	return math.Hypot(real(c), imag(c))
}
