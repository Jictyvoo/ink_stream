package imageparser

import (
	"image"

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
	}
)

func NewImagePipeline(steps ...PipeStep) ImagePipeline {
	totalSteps := uint8(len(steps))
	imgPipe := ImagePipeline{
		fullProcessSteps: make([]PipeStep, 0, totalSteps>>1),
		pixelSteps:       make([]UnitStep, 0, totalSteps>>1),
	}

	for _, step := range steps {
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
		newImage := createDrawImage(state.img, img.Bounds())
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
