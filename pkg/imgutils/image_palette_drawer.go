package imgutils

import (
	"image"
	"image/color"
	"image/draw"
)

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

func NewDrawFromImgColorModel(colorModel color.Model, bounds image.Rectangle) draw.Image {
	switch colorModel {
	case color.GrayModel, color.Gray16Model:
		return image.NewGray(bounds)
	case color.RGBAModel, color.RGBA64Model, color.NRGBAModel, color.NRGBA64Model:
		return image.NewRGBA(bounds)
	}

	return image.NewRGBA(bounds)
}

func (fac imageFactory) CreateDrawImage(colorModel color.Model, bounds image.Rectangle) draw.Image {
	newImg := NewDrawFromImgColorModel(colorModel, bounds)
	return ImagePaletteDrawer{palette: fac.palette, Image: newImg}
}

func (i ImagePaletteDrawer) Set(x, y int, c color.Color) {
	newColor := c
	if len(i.palette) > 0 {
		newColor = i.palette.Convert(c)
	}

	i.Image.Set(x, y, newColor)
}
