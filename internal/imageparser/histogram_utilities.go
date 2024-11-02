package imageparser

import "image"

const maxPixelValue = (1 << 8) - 1

type (
	ChannelHistogram [maxPixelValue + 1]uint32
	ImageHistogram   struct {
		data [3]ChannelHistogram
	}
)

// calculateHistogram calculates the histogram for each color channel in an image.
func calculateHistogram(img image.Image) ImageHistogram {
	var histogram [3]ChannelHistogram
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			histogram[0][r>>8]++
			histogram[1][g>>8]++
			histogram[2][b>>8]++
		}
	}

	return ImageHistogram{data: histogram}
}

func (histogram ImageHistogram) hiloHistogram(
	minVal [3]uint8, maxVal [3]uint8,
	stopChannels [3]struct{ min, max bool },
) ([3]uint8, [3]uint8) {
	for (!stopChannels[0].min || !stopChannels[1].min || !stopChannels[2].min) ||
		(!stopChannels[0].max || !stopChannels[1].max || !stopChannels[2].max) {
		for index := range 3 {
			histogram.channelHiLo(
				histogram.data[index],
				&minVal[index], &maxVal[index],
				&stopChannels[index],
			)
		}
	}

	return minVal, maxVal
}

func (histogram ImageHistogram) channelHiLo(
	channelData ChannelHistogram, minVal, maxVal *uint8,
	stop *struct{ min, max bool },
) {
	if !stop.min {
		stop.min = channelData[*minVal] > 0
		if !stop.min { // Check again to ensure value is correct
			*minVal++
			stop.min = *minVal == maxPixelValue
		}
	}

	if !stop.max {
		stop.max = channelData[*maxVal] > 0
		if !stop.max { // Check again to ensure value is correct
			*maxVal--
			stop.max = *maxVal == maxPixelValue
		}
	}
}
