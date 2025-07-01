package testimgs

import (
	"crypto/sha256"
	"image"
	"image/color"
	"image/draw"
	"math/rand/v2"
)

func ImageFixtures(total uint8, seed []byte) []image.Image {
	resultImgs := make([]image.Image, 0, total)

	// Seed the random number generator with the hash of the seed to ensure sufficient entropy.
	hash := sha256.Sum256(seed)
	rng := rand.New(rand.NewChaCha8(hash))

	for i := uint8(0); i < total; i++ {
		// Random width and height between 8 and 32 pixels
		width := 8 + rng.IntN(25)  // min 8, max 32
		height := 8 + rng.IntN(25) // min 8, max 32

		// Create a blank RGBA image with the random dimensions
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// Fill the image with random colors
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := uint8(rng.IntN(256))
				g := uint8(rng.IntN(256))
				b := uint8(rng.IntN(256))
				a := uint8(rng.IntN(256))
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
			}
		}

		// Add to the slice of resultImgs
		resultImgs = append(resultImgs, img)
	}

	return resultImgs
}

// NewSolidImage creates a new image with the specified bounds and fills it with the given solid color.
func NewSolidImage(bounds image.Rectangle, fillColor color.Color) image.Image {
	img := image.NewRGBA(bounds)
	col := image.NewUniform(fillColor)

	// Draw the uniform color over the entire image
	draw.Draw(img, bounds, col, image.Point{}, draw.Src)
	return img
}

// NewBorderedImage creates an image filled with a background color,
// and a central content area (inside the given margins) filled with a foreground color.
func NewBorderedImage(
	imgRect image.Rectangle,
	top, bottom, left, right int,
	bg, fg color.Color,
) image.Image {
	img := image.NewNRGBA(imgRect)

	// Fill entire image with background color
	draw.Draw(img, imgRect, image.NewUniform(bg), image.Point{}, draw.Src)

	// Define the content rectangle based on margins
	contentRect := image.Rect(
		imgRect.Min.X+left,
		imgRect.Min.Y+top,
		imgRect.Max.X-right,
		imgRect.Max.Y-bottom,
	)

	// Fill content area with foreground color
	draw.Draw(img, contentRect, image.NewUniform(fg), image.Point{}, draw.Src)

	return img
}
