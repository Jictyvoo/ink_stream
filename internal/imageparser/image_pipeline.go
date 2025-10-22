package imageparser

import (
	"image"
	"image/color"
	"slices"

	"github.com/Jictyvoo/ink_stream/pkg/imgutils"
)

type (
	PipeState struct {
		name      string
		Img       image.Image
		SubImages []image.Image
	}
	ProcessOptions struct {
		Gamma float64
	}
	ImagePipeline struct {
		opts             ProcessOptions
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

func (imgPipe ImagePipeline) processImage(
	img image.Image, skipSteps []string,
) (resultImg image.Image, subImages []image.Image, executedSteps []string, err error) {
	state := PipeState{Img: img}
	executedSteps = make([]string, 0, len(imgPipe.fullProcessSteps))
	for _, step := range imgPipe.fullProcessSteps {
		if slices.Contains(skipSteps, step.StepID()) {
			continue
		}
		if err = step.PerformExec(&state, imgPipe.opts); err != nil {
			return resultImg, subImages, executedSteps, err
		}
		executedSteps = append(executedSteps, step.StepID())
		if len(state.SubImages) > 0 {
			subImages = append(subImages, state.SubImages...)
			return resultImg, subImages, executedSteps, err
		}
	}

	// Check if it has pixel steps
	if len(imgPipe.pixelSteps) > 0 {
		img = state.Img
		newImage := imgPipe.drawFactory.CreateDrawImage(state.Img.ColorModel(), img.Bounds())
		for x, y := range imgutils.Iterator(img) {
			oldColor := img.At(x, y)
			for _, step := range imgPipe.pixelSteps {
				oldColor = step.PixelStep(oldColor)
			}
			newImage.Set(x, y, oldColor)
		}

		state.Img = newImage
	}

	resultImg = state.Img
	return resultImg, subImages, executedSteps, err
}

func (imgPipe ImagePipeline) Process(img image.Image) (outputImgs []image.Image, err error) {
	imgSlice := []image.Image{img}
	var skipSteps []string
	for index := 0; index < len(imgSlice); index++ {
		singleImage, subImages, executedSteps, processErr := imgPipe.processImage(
			imgSlice[index], skipSteps,
		)
		if processErr != nil {
			return nil, err
		}

		if singleImage != nil {
			outputImgs = append(outputImgs, singleImage)
		}

		if index >= len(imgSlice)-1 {
			clear(imgSlice) // Try to free the memory
		}

		if len(subImages) > 0 {
			imgSlice = append(imgSlice, subImages...)
			skipSteps = append(skipSteps, executedSteps...)
		}
	}

	return outputImgs, err
}

func (imgPipe ImagePipeline) PipeSteps() []PipeStep {
	return imgPipe.fullProcessSteps
}
