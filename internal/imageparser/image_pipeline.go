package imageparser

import "image"

type (
	pipeState struct {
		name string
		img  image.Image
	}
	processOptions struct {
		gamma      float64
		applyColor bool
	}
	PipeStep interface {
		PerformExec(state *pipeState, opts processOptions) (err error)
	}
	ImagePipeline struct {
		state pipeState
		opts  processOptions
		steps []PipeStep
	}
)

func (imgPipe ImagePipeline) Process(img image.Image) (resultImg image.Image, err error) {
	imgPipe.state.img = img
	for _, step := range imgPipe.steps {
		if err = step.PerformExec(&imgPipe.state, imgPipe.opts); err != nil {
			return
		}
	}

	resultImg = imgPipe.state.img
	return
}
