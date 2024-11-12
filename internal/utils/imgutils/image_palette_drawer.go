package imgutils

import (
	"image"
	"image/color"
	"image/draw"
)

type DrawImageFactory interface {
	CreateDrawImage(img image.Image, bounds image.Rectangle) draw.Image
}

type (
	ImagePaletteDrawer struct {
		palette color.Palette
		draw.Image
	}
	imageFactory struct {
		palette color.Palette
	}
)

func NewImageFactory(palette color.Palette) DrawImageFactory {
	return &imageFactory{palette: palette}
}

func (fac imageFactory) drawImageFromImage(img image.Image, bounds image.Rectangle) draw.Image {
	switch img.ColorModel() {
	case color.GrayModel, color.Gray16Model:
		return image.NewGray(bounds)
	case color.RGBAModel, color.RGBA64Model, color.NRGBAModel, color.NRGBA64Model:
		return image.NewRGBA(bounds)
	}

	return image.NewRGBA(bounds)
}

func (fac imageFactory) CreateDrawImage(img image.Image, bounds image.Rectangle) draw.Image {
	newImg := fac.drawImageFromImage(img, bounds)
	return ImagePaletteDrawer{palette: fac.palette, Image: newImg}
}

func (i ImagePaletteDrawer) Set(x, y int, c color.Color) {
	newColor := c
	if len(i.palette) > 0 {
		newColor = i.palette.Convert(c)
	}

	i.Image.Set(x, y, newColor)
}
