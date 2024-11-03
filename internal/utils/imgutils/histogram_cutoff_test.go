package imgutils

import (
	"testing"
)

func TestApplyCutoff(t *testing.T) {
	tests := []struct {
		name           string
		histogram      ChannelHistogram
		cutoffLow      float64
		cutoffHigh     float64
		expectedResult ChannelHistogram
	}{
		{
			name: "No cutoff applied",
			histogram: ChannelHistogram{
				0: 10, 1: 20, 255: 15, // only a few values set, the rest are zero
			},
			cutoffLow:      0,
			cutoffHigh:     0,
			expectedResult: ChannelHistogram{0: 10, 1: 20, 255: 15},
		},
		{
			name: "Apply low cutoff",
			histogram: ChannelHistogram{
				0: 50, 1: 100, 2: 150, 3: 200, 4: 500, 5: 1000, 255: 300,
			},
			cutoffLow:  10, // 10% of the total pixels
			cutoffHigh: 0,  // no high cutoff
			expectedResult: ChannelHistogram{
				0: 0, 1: 0, 2: 70, 3: 200, 4: 500, 5: 1000, 255: 300,
			}, // Low values should be partially zeroed
		},
		{
			name: "Apply high cutoff",
			histogram: ChannelHistogram{
				0: 50, 1: 100, 2: 150, 3: 200, 4: 500, 5: 1000, 255: 300,
			},
			cutoffLow:  0,  // no low cutoff
			cutoffHigh: 10, // 10% of the total pixels
			expectedResult: ChannelHistogram{
				0: 50, 1: 100, 2: 150, 3: 200, 4: 500, 5: 1000, 255: 70,
			}, // High values should be partially zeroed
		},
		{
			name: "Apply both low and high cutoffs",
			histogram: ChannelHistogram{
				0: 50, 1: 100, 2: 150, 3: 200, 4: 500, 5: 1000, 6: 700, 255: 300,
			},
			cutoffLow:  5,  // 5% of the total pixels
			cutoffHigh: 10, // 10% of the total pixels
			expectedResult: ChannelHistogram{
				0: 0, 1: 0, 2: 150, 3: 200, 4: 500, 5: 1000, 6: 700, 255: 0,
			}, // Low and high values should be partially zeroed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyCutoff(tt.histogram, tt.cutoffLow, tt.cutoffHigh)

			for i, v := range tt.expectedResult {
				if result[i] != v {
					t.Errorf("expected histogram[%d] to be %d, got %d", i, v, result[i])
				}
			}
		})
	}
}
