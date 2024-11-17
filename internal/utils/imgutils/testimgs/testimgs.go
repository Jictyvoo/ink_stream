package testimgs

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
)

var (
	//go:embed black_square_gray_left_green_right_white_margin_right.png
	blackSquareGrayLeftGreenRightWhiteMargin []byte

	//go:embed black_square_right_green_gray_margins.png
	blackSquareRightGreenGrayMargins []byte
)

func ImageBlackSquareGreenRight(whiteMargin bool) image.Image {
	imgBytes := blackSquareRightGreenGrayMargins
	if whiteMargin {
		imgBytes = blackSquareGrayLeftGreenRightWhiteMargin
	}

	img, _ := png.Decode(bytes.NewReader(imgBytes))
	return img
}

//go:embed black_square_white_margins.png
var blackSquareWhiteMargins []byte

func ImageBlackSquareWhiteMargin() image.Image {
	img, _ := png.Decode(bytes.NewReader(blackSquareWhiteMargins))
	return img
}

//go:embed black_circle_with_transparent_background.png
var blackCircleWithTransparentBackground []byte

func ImageBlackCircleWithTransparentBackground() image.Image {
	img, _ := png.Decode(bytes.NewReader(blackCircleWithTransparentBackground))
	return img
}
