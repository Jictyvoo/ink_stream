package testimgs

import (
	"bytes"
	"image"
	"testing"
	"time"
)

// Test that the same seed generates the same images
func TestImageFixtures_Reproducibility(t *testing.T) {
	seed := []byte("reproducible-seed")
	total := uint8(3)

	images1 := ImageFixtures(total, seed)
	images2 := ImageFixtures(total, seed)

	for i := range images1 {
		img1 := images1[i].(*image.RGBA)
		img2 := images2[i].(*image.RGBA)

		if !bytes.Equal(img1.Pix, img2.Pix) {
			t.Errorf("Images at index %d do not match for the same seed", i)
		}
	}
}

// Test that different seeds produce different images
func TestImageFixtures_DifferentSeeds(t *testing.T) {
	total := uint8(3)
	seed1 := []byte("seed-one")
	seed2 := []byte("seed-two")

	images1 := ImageFixtures(total, seed1)
	images2 := ImageFixtures(total, seed2)

	different := false
	for i := range images1 {
		img1 := images1[i].(*image.RGBA)
		img2 := images2[i].(*image.RGBA)

		if !bytes.Equal(img1.Pix, img2.Pix) {
			different = true
			break
		}
	}

	if !different {
		t.Errorf("Expected different images for different seeds")
	}
}

func TestImageFixtures(t *testing.T) {
	seed := []byte("test-seed")
	total := uint8(5)

	// Generate images with the ImageFixtures function
	originalImgSet := ImageFixtures(total, seed)

	// 1. Test the number of generated originalImgSet
	if len(originalImgSet) != int(total) {
		t.Fatalf("Expected %d originalImgSet, got %d", total, len(originalImgSet))
	}

	time.Sleep(time.Millisecond << 4)

	regeneratedImgSet := ImageFixtures(total, seed)
	seed2 := []byte("different-seed")
	otherSeedImgSet := ImageFixtures(total, seed2)
	different := false

	for i, img := range originalImgSet {
		// 2. Test dimensions of each generated image
		{
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			if width < 8 || width > 32 || height < 8 || height > 32 {
				t.Errorf("Image %d has invalid dimensions: width=%d, height=%d", i, width, height)
			}
		}

		// 3. Test pixel values are within range [0, 255]
		{
			rgbaImg, ok := img.(*image.RGBA)
			if !ok {
				t.Fatalf("Image %d is not of type RGBA", i)
			}
			for _, v := range rgbaImg.Pix {
				if v < 0 || v > 255 {
					t.Errorf("Image %d has pixel value out of range: %d", i, v)
				}
			}
		}

		// 4. Test reproducibility with the same seed
		{
			img1 := originalImgSet[i].(*image.RGBA)
			img2 := regeneratedImgSet[i].(*image.RGBA)

			if !bytes.Equal(img1.Pix, img2.Pix) {
				t.Errorf("Images at index %d do not match for the same seed", i)
			}
		}

		// 5. Test that different seeds produce different originalImgSet
		{
			img1 := originalImgSet[i].(*image.RGBA)
			img3 := otherSeedImgSet[i].(*image.RGBA)

			if !bytes.Equal(img1.Pix, img3.Pix) {
				different = true
				break
			}
		}
	}

	if !different {
		t.Errorf("Expected different originalImgSet for different seeds")
	}
}
