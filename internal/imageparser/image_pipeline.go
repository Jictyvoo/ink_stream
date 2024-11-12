package imageparser

import (
	"image"
	"image/color"

	"github.com/Jictyvoo/ink_stream/internal/utils/imgutils"
)

type (
	pipeState struct {
		name string
		img  image.Image
	}
	processOptions struct {
		gamma      float64
		applyColor bool
	}
	ImagePipeline struct {
		opts             processOptions
		pixelSteps       []UnitStep
		fullProcessSteps []PipeStep
		drawFactory      imgutils.DrawImageFactory
	}
)

func NewImagePipeline(palette color.Palette, steps ...PipeStep) ImagePipeline {
	imgPipe := ImagePipeline{
		fullProcessSteps: steps,
		pixelSteps:       []UnitStep{},
		drawFactory:      imgutils.NewImageFactory(palette),
	}

	for _, step := range imgPipe.fullProcessSteps {
		step.UpdateDrawFactory(imgPipe.drawFactory)
	}

	return imgPipe
}

func NewImagePipelineSplitStep(palette color.Palette, steps ...PipeStep) ImagePipeline {
	totalSteps := uint8(len(steps))
	imgPipe := ImagePipeline{
		fullProcessSteps: make([]PipeStep, 0, totalSteps>>1),
		pixelSteps:       make([]UnitStep, 0, totalSteps>>1),
		drawFactory:      imgutils.NewImageFactory(palette),
	}

	for _, step := range steps {
		step.UpdateDrawFactory(imgPipe.drawFactory)
		switch objType := step.(type) {
		case UnitStep:
			imgPipe.pixelSteps = append(imgPipe.pixelSteps, objType)
		default:
			imgPipe.fullProcessSteps = append(imgPipe.fullProcessSteps, step)
		}
	}

	return imgPipe
}

func (imgPipe ImagePipeline) Process(img image.Image) (resultImg image.Image, err error) {
	state := pipeState{img: img}
	for _, step := range imgPipe.fullProcessSteps {
		if err = step.PerformExec(&state, imgPipe.opts); err != nil {
			return
		}
	}

	// Check if it has pixel steps
	if len(imgPipe.pixelSteps) > 0 {
		img = state.img
		newImage := imgPipe.drawFactory.CreateDrawImage(state.img, img.Bounds())
		for x, y := range imgutils.Iterator(img) {
			oldColor := img.At(x, y)
			for _, step := range imgPipe.pixelSteps {
				oldColor = step.PixelStep(oldColor)
			}
			newImage.Set(x, y, oldColor)
		}

		state.img = newImage
	}

	resultImg = state.img
	return
}
