package imgutils

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/utils"
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
	histogram := CalculateHistogram(img)

	// Verify the red channel has counts at index 255
	if histogram.data[0][255] != 4 {
		t.Errorf("expected 4 red pixels at 255, got %d", histogram.data[0][255])
	}

	// Verify the green and blue channels are all zero
	for i := 0; i <= MaxPixelValue; i++ {
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
		stop   utils.MinMaxGeneric[bool]
	)

	// Test finding the min and max values in the histogram
	for range MaxPixelValue {
		channelHiLo(histogram, &minVal, &maxVal, &stop)
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
		minVal [3]uint8
		maxVal = [3]uint8{255, 255, 255}
	)

	// Calculate hilo histogram
	minResult, maxResult := histogram.HiloHistogram(minVal, maxVal)

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

func TestImageHistogram_Channel(t *testing.T) {
	// Arrange
	hist := ImageHistogram{
		data: [3]ChannelHistogram{
			{1, 2, 3}, {4, 5, 6}, {7, 8, 9},
		},
	}

	tests := []struct {
		name      string
		index     uint8
		expected  ChannelHistogram
		expectNil bool
	}{
		{
			name:     "Valid index 0",
			index:    0,
			expected: hist.data[0],
		},
		{
			name:     "Valid index 1",
			index:    1,
			expected: hist.data[1],
		},
		{
			name:     "Valid index 2",
			index:    2,
			expected: hist.data[2],
		},
		{
			name:      "Index out of range",
			index:     3,
			expectNil: true,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hist.Channel(tt.index)

			if tt.expectNil {
				if !slices.Equal(result[:], make([]uint32, 256)) {
					t.Errorf(
						"Expected empty ChannelHistogram for out-of-range index %d, got %v",
						tt.index, result,
					)
				}
				return
			}

			if !slices.Equal(result[:], tt.expected[:]) {
				t.Errorf(
					"Expected ChannelHistogram for index %d, got %v",
					tt.index, result,
				)
			}
		})
	}
}

func TestImageHistogram_Set(t *testing.T) {
	newChannel := ChannelHistogram{1 << 4, 3 << 9, 7 << 1}

	tests := []struct {
		name     string
		index    uint8
		channel  ChannelHistogram
		expected []ChannelHistogram
	}{
		{
			name:     "Set within range",
			index:    0,
			channel:  newChannel,
			expected: []ChannelHistogram{newChannel, {4, 5, 6}, {7, 8, 9}},
		},
		{
			name:     "Another set within range",
			index:    1,
			channel:  newChannel,
			expected: []ChannelHistogram{{1, 2, 3}, newChannel, {7, 8, 9}},
		},
		{
			name:     "Index out of range",
			index:    3,
			channel:  newChannel,
			expected: []ChannelHistogram{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, // should be unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := ImageHistogram{
				data: [3]ChannelHistogram{
					{1, 2, 3}, {4, 5, 6}, {7, 8, 9},
				},
			}
			hist.Set(tt.index, tt.channel)

			// Loop through each expected ChannelHistogram and assert equality
			for i, expectedChannel := range tt.expected {
				if hist.data[i] != expectedChannel {
					t.Errorf(
						"Test %q failed at index %d: expected %v, got %v",
						tt.name, i,
						expectedChannel, hist.data[i],
					)
				}
			}
		})
	}
}
