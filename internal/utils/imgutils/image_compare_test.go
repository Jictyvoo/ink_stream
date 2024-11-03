package imgutils

import (
	_ "embed"
	"image"
	"image/color"
	"slices"
	"testing"
)

//go:embed image_compare_test.go
var imageCompareFile []byte

func TestIsImageEqual(t *testing.T) {
	const totalFixtures = 6
	imgFixtures := ([totalFixtures]image.Image)(ImageFixtures(totalFixtures, imageCompareFile))
	var whiteImagePix [24]uint8
	for i := range whiteImagePix {
		whiteImagePix[i] = 255
	}

	testCases := []struct {
		name     string
		a, b     image.Image
		expected bool
	}{
		{
			name:     "Same images",
			a:        imgFixtures[5],
			b:        imgFixtures[5],
			expected: true,
		},
		{
			name:     "Identical empty RGBA images",
			a:        &image.RGBA{},
			b:        &image.RGBA{},
			expected: true,
		},
		{
			name: "Empty RGBA and random image comparison",
			a:    &image.RGBA{},
			b:    imgFixtures[4],
		},
		{
			name: "1x1 black pixel vs different color",
			a: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
				newImg.Set(0, 0, color.Black)
				return newImg
			}(),
			b: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
				newImg.Set(0, 0, color.RGBA{R: 1, A: 255})
				return newImg
			}(),
			expected: false,
		},
		{
			name: "Identical 3x2 white images",
			a: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 3, 2))
				newImg.Pix = whiteImagePix[:]
				return newImg
			}(),
			b: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 3, 2))
				newImg.Pix = slices.Clone(whiteImagePix[:])
				newImg.Pix[1] = 0
				return newImg
			}(),
			expected: false,
		},
		{
			name: "Different sizes (3x2 and 2x2)",
			a: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 3, 2))
				newImg.Pix = whiteImagePix[:]
				return newImg
			}(),
			b: func() image.Image {
				newImg := image.NewRGBA(image.Rect(0, 0, 2, 2))
				newImg.Pix = whiteImagePix[:16]
				return newImg
			}(),
			expected: false,
		},
		{
			name:     "Fixture images 0 and 1",
			a:        imgFixtures[0],
			b:        imgFixtures[1],
			expected: false, // Assuming different fixture images are generated
		},
		{
			name:     "Fixture images 2 and 3",
			a:        imgFixtures[2],
			b:        imgFixtures[3],
			expected: false, // Assuming different fixture images are generated
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actual := IsImageEqual(tCase.a, tCase.b)
			if actual != tCase.expected {
				t.Errorf("%s: expected: %v actual: %v", "IsImageEqual", tCase.expected, actual)
			}
		})
	}
}

func TestIterator_CompleteIteration(t *testing.T) {
	const width, height = 3, 3
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Create a map to keep track of visited coordinates
	visited := make(map[uint16]bool)
	visitIndex := func(x, y int) uint16 {
		return (uint16(x) << 8) | uint16(y)
	}

	// Run the iterator
	Iterator(img)(func(x, y int) bool {
		visited[visitIndex(x, y)] = true
		return true
	})

	// Check if all coordinates within bounds were visited
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !visited[visitIndex(x, y)] {
				t.Errorf("Coordinate (%d, %d) was not visited", x, y)
			}
		}
	}
}

func TestIterator_EarlyExit(t *testing.T) {
	const width, height = 3, 3
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Keep track of number of visited coordinates
	count := 0
	expectedVisits := 1 // We expect only one visit since we return false immediately

	// Run the iterator with early exit
	Iterator(img)(func(x, y int) bool {
		count++
		return false // Exit after first coordinate
	})

	// Check that only one coordinate was visited
	if count != expectedVisits {
		t.Errorf("Expected %d coordinates to be visited, but got %d", expectedVisits, count)
	}
}
