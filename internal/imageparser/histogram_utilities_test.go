package imageparser

import (
	"image"
	"image/color"
	"testing"
)

// Helper function to create a solid color image for testing.
func createSolidColorImage(width, height int, col color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, col)
		}
	}
	return img
}

// TestCalculateHistogram tests the calculateHistogram function.
func TestCalculateHistogram(t *testing.T) {
	// Create a small test image with a solid red color
	img := createSolidColorImage(2, 2, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Calculate histogram
	histogram := calculateHistogram(img)

	// Verify the red channel has counts at index 255
	if histogram.data[0][255] != 4 {
		t.Errorf("expected 4 red pixels at 255, got %d", histogram.data[0][255])
	}

	// Verify the green and blue channels are all zero
	for i := 0; i <= maxPixelValue; i++ {
		if i == 0 {
			var zeroes struct{ green, blue uint32 }
			zeroes.green, zeroes.blue = histogram.data[1][i], histogram.data[2][i]
			if zeroes.green == 0 {
				t.Errorf("expected green pixel at 0, got %d", zeroes.green)
			}
			if zeroes.blue == 0 {
				t.Errorf("expected blue pixel at 0, got %d", zeroes.blue)
			}
			continue
		}
		if histogram.data[1][i] != 0 {
			t.Errorf("expected 0 green pixels, got %d at index %d", histogram.data[1][i], i)
		}
		if histogram.data[2][i] != 0 {
			t.Errorf("expected 0 blue pixels, got %d at index %d", histogram.data[2][i], i)
		}
	}
}

// TestChannelHiLo tests the channelHiLo function.
func TestChannelHiLo(t *testing.T) {
	histogram := ChannelHistogram{
		10: 5, 100: 3,
		200: 10, 255: 0,
	}

	var (
		minVal uint8 = 0
		maxVal uint8 = 255
		stop   struct{ min, max bool }
	)

	// Test finding the min and max values in the histogram
	imgHist := ImageHistogram{}
	for range maxPixelValue {
		imgHist.channelHiLo(histogram, &minVal, &maxVal, &stop)
	}

	if minVal != 10 {
		t.Errorf("expected minVal 10, got %d", minVal)
	}
	if maxVal != 200 {
		t.Errorf("expected maxVal 200, got %d", maxVal)
	}
}

// TestHiloHistogram tests the hiloHistogram function.
func TestHiloHistogram(t *testing.T) {
	histogram := ImageHistogram{
		data: [3]ChannelHistogram{
			{0: 0, 86: 1, 121: 1, 64: 1},   // Red channel
			{10: 3, 50: 1, 150: 1, 189: 1}, // Green channel
			{0: 0, 50: 1, 100: 1, 201: 1},  // Blue channel
		},
	}

	// Define initial min and max values and stop channels
	var (
		minVal       [3]uint8
		maxVal       = [3]uint8{255, 255, 255}
		stopChannels [3]struct{ min, max bool }
	)

	// Calculate hilo histogram
	minResult, maxResult := histogram.hiloHistogram(minVal, maxVal, stopChannels)

	// Expected min and max values based on histogram data
	expectedMin := [3]uint8{64, 10, 50}
	expectedMax := [3]uint8{121, 189, 201}

	// Check if the results match expected values
	for i := 0; i < 3; i++ {
		if minResult[i] != expectedMin[i] {
			t.Errorf(
				"expected min value for channel %d: %d, got %d",
				i, expectedMin[i], minResult[i],
			)
		}
		if maxResult[i] != expectedMax[i] {
			t.Errorf(
				"expected max value for channel %d: %d, got %d",
				i, expectedMax[i], maxResult[i],
			)
		}
	}
}
