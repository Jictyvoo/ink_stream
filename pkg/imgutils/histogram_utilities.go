package imgutils

import (
	"image"

	"github.com/Jictyvoo/ink_stream/internal/utils"
)

type (
	ChannelHistogram [MaxPixelValue + 1]uint32
	ImageHistogram   struct {
		data [3]ChannelHistogram
	}
)

// CalculateHistogram calculates the histogram for each color channel in an image.
func CalculateHistogram(img image.Image) ImageHistogram {
	var histogram [3]ChannelHistogram

	for x, y := range Iterator(img) {
		r, g, b, _ := img.At(x, y).RGBA()
		histogram[0][r>>8]++
		histogram[1][g>>8]++
		histogram[2][b>>8]++
	}

	return ImageHistogram{data: histogram}
}

func (histogram ImageHistogram) HiloHistogram() (minVal [3]uint8, maxVal [3]uint8) {
	maxVal = [3]uint8{MaxPixelValue, MaxPixelValue, MaxPixelValue}
	var stopChannels [3]utils.MinMaxGeneric[bool]
	for (!stopChannels[0].Min || !stopChannels[1].Min || !stopChannels[2].Min) ||
		(!stopChannels[0].Max || !stopChannels[1].Max || !stopChannels[2].Max) {
		for index := range 3 {
			channelHiLo(
				histogram.data[index],
				&minVal[index], &maxVal[index],
				&stopChannels[index],
			)
		}
	}

	return minVal, maxVal
}

func (histogram *ImageHistogram) Channel(index uint8) ChannelHistogram {
	if index >= uint8(len(histogram.data)) {
		return ChannelHistogram{}
	}

	return histogram.data[index]
}

func (histogram *ImageHistogram) Set(index uint8, channel ChannelHistogram) bool {
	if index >= uint8(len(histogram.data)) {
		return false
	}

	histogram.data[index] = channel
	return true
}

func channelHiLo(
	channelData ChannelHistogram,
	minVal, maxVal *uint8,
	stop *utils.MinMaxGeneric[bool],
) {
	if !stop.Min {
		stop.Min = channelData[*minVal] > 0
		if !stop.Min { // Check again to ensure value is correct
			*minVal++
			stop.Min = *minVal == MaxPixelValue
		}
	}

	if !stop.Max {
		stop.Max = channelData[*maxVal] > 0
		if !stop.Max { // Check again to ensure value is correct
			*maxVal--
			stop.Max = *maxVal == MaxPixelValue
		}
	}
}
