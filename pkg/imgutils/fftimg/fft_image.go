package fftimg

import (
	"image"
	"image/color"
	"image/draw"
	"sync/atomic"

	"gonum.org/v1/gonum/dsp/fourier"
)

// FFTAccessible defines access to the FFT of an image.
type FFTAccessible interface {
	FFT() [4][][]complex128 // [channel][y][x]
}

// FFTImage wraps a draw.Image and lazily computes its FFT per channel.
type FFTImage struct {
	img      draw.Image
	fftData  [4][][]complex128 // R, G, B, A
	fftDirty atomic.Bool       // true if FFT must be recalculated
	w, h     int
}

// NewFFTImage creates a new FFTImage from a draw.Image and pre-allocates fftData.
func NewFFTImage(img draw.Image) *FFTImage {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	f := &FFTImage{
		img: img,
		w:   w,
		h:   h,
	}
	for i := 0; i < 4; i++ {
		f.fftData[i] = make([][]complex128, h)
		for y := 0; y < h; y++ {
			f.fftData[i][y] = make([]complex128, w)
		}
	}
	f.fftDirty.Store(true)
	return f
}

// ColorModel implements draw.Image.
func (f *FFTImage) ColorModel() color.Model {
	return f.img.ColorModel()
}

// Bounds implements draw.Image.
func (f *FFTImage) Bounds() image.Rectangle {
	return f.img.Bounds()
}

// At implements draw.Image.
func (f *FFTImage) At(x, y int) color.Color {
	return f.img.At(x, y)
}

// Set implements draw.Image. Marks FFT as dirty.
func (f *FFTImage) Set(x, y int, c color.Color) {
	f.img.Set(x, y, c)
	f.fftDirty.Store(true)
}

// FFT returns the 2D FFT of the image for each channel, computing it lazily.
func (f *FFTImage) FFT() [4][][]complex128 {
	if f.fftDirty.Load() {
		f.computeImageFFT2()
		f.fftDirty.Store(false)
	}
	return f.fftData
}

// computeImageFFT2 computes the 2D FFT for each channel (R, G, B, A) using gonum.org/v1/gonum/dsp/fourier.
func (f *FFTImage) computeImageFFT2() {
	b := f.img.Bounds()
	w, h := f.w, f.h
	// Prepare per-channel float64 matrices
	channels := [4][][]float64{}
	for i := 0; i < 4; i++ {
		channels[i] = make([][]float64, h)
		for y := 0; y < h; y++ {
			channels[i][y] = make([]float64, w)
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := f.img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			channels[0][y][x] = float64(r >> 8)
			channels[1][y][x] = float64(g >> 8)
			channels[2][y][x] = float64(b >> 8)
			channels[3][y][x] = float64(a >> 8)
		}
	}
	// For each channel, perform 2D FFT
	for ch := 0; ch < 4; ch++ {
		// First axis (rows)
		fft := fourier.NewFFT(w)
		rowCoeffs := make([][]complex128, h)
		for y := 0; y < h; y++ {
			rowCoeffs[y] = make([]complex128, w)
			fft.Coefficients(rowCoeffs[y], channels[ch][y])
		}
		// Second axis (columns)
		cfft := fourier.NewCmplxFFT(h)
		col := make([]complex128, h)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				col[y] = rowCoeffs[y][x]
			}
			cfft.Coefficients(col, col)
			for y := 0; y < h; y++ {
				f.fftData[ch][y][x] = col[y]
			}
		}
	}
}
