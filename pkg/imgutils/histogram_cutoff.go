package imgutils

// ApplyCutoff modifies the histogram by cutting off a percentage of pixels at both low and high ends.
func ApplyCutoff(histogram ChannelHistogram, cutoffLow, cutoffHigh float64) ChannelHistogram {
	// Calculate the total number of pixels
	var n uint32
	for _, count := range histogram {
		n += count
	}

	// Apply cutoff for the low end
	var (
		cut struct{ high, low uint32 }
		lo  uint32 = 0
		hi  uint32 = MaxPixelValue
	)
	cut.low = uint32(float64(n)*cutoffLow) / 100
	cut.high = uint32(float64(n)*cutoffHigh) / 100

	var (
		controlStop [2]bool
		cutFunc     = func(iterVal uint32, cutValue *uint32, shouldStop *bool) {
			if !*shouldStop {
				if *cutValue > histogram[iterVal] {
					*cutValue -= histogram[iterVal]
					histogram[iterVal] = 0
				} else {
					histogram[iterVal] -= *cutValue
					*cutValue = 0
				}
				*shouldStop = *cutValue <= 0
			}
		}
	)
	for range len(histogram) {
		// Apply cutoff for the low end
		cutFunc(lo, &cut.low, &controlStop[0])

		// Apply cutoff for the high end
		cutFunc(hi, &cut.high, &controlStop[1])

		if controlStop[0] && controlStop[1] {
			break
		}
		lo++
		hi--
	}

	return histogram
}
