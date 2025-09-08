package inktypes

import (
	"image/color"
	"testing"
)

func TestPaletteIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		palette     color.Palette
		expectedHex string
	}{
		{
			name:        "empty palette",
			palette:     color.Palette{},
			expectedHex: "",
		},
		{
			name: "single red",
			palette: color.Palette{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
			expectedHex: "6666303030306666",
		},
		{
			name: "red and green",
			palette: color.Palette{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 0, A: 255},
			},
			expectedHex: "66663030303066662d6666303066663030",
		},
		{
			name: "blue with half alpha",
			palette: color.Palette{
				color.RGBA{R: 0, G: 0, B: 255, A: 128},
			},
			expectedHex: "3830666630303030",
		},
		{
			name: "three colors",
			palette: color.Palette{
				color.RGBA{R: 255, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 255, A: 255},
				color.RGBA{R: 255, G: 0, B: 255, A: 255},
			},
			expectedHex: "66663030666666662d66666666666630302d6666666630306666",
		},
		{
			name: "duplicate colors",
			palette: color.Palette{
				color.RGBA{R: 10, G: 20, B: 30, A: 40},
				color.RGBA{R: 10, G: 20, B: 30, A: 40},
			},
			expectedHex: "32383165313430612d3238316531343061",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := NewPaletteIdentifier(tc.palette)

			if gotHex := id.Hex(); gotHex != tc.expectedHex {
				t.Fatalf("Hex() = %q, want %q", gotHex, tc.expectedHex)
			}
		})
	}
}

func TestPaletteRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		palette color.Palette
	}{
		{
			name:    "Single color opaque",
			palette: color.Palette{color.NRGBA{R: 255, G: 0, B: 0, A: 255}}, // red
		},
		{
			name: "Two colors",
			palette: color.Palette{
				color.NRGBA{R: 0, G: 255, B: 0, A: 255}, // green
				color.NRGBA{R: 0, G: 0, B: 255, A: 255}, // blue
			},
		},
		{
			name: "With transparency",
			palette: color.Palette{
				color.NRGBA{R: 123, G: 45, B: 67, A: 128}, // semi-transparent
				color.NRGBA{R: 200, G: 100, B: 50, A: 0},  // fully transparent
			},
		},
		{
			name:    "Empty palette",
			palette: color.Palette{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			pi := NewPaletteIdentifier(tt.palette)

			// Decode
			got, err := pi.ToPalette()
			if err != nil {
				t.Fatalf("ToPalette() error = %v", err)
			}

			if len(got) != len(tt.palette) {
				t.Fatalf(
					"Round-trip mismatch. Got %d colors, expected %d",
					len(got),
					len(tt.palette),
				)
			}

			for index, expectedColor := range tt.palette {
				var compareColors struct {
					expected, received struct{ r, g, b, a uint32 }
				}

				compareColors.expected.r, compareColors.expected.g, compareColors.expected.b, compareColors.expected.a = expectedColor.RGBA()
				compareColors.received.r, compareColors.received.g, compareColors.received.b, compareColors.received.a = got[index].RGBA()
				shiftPixel(
					&compareColors.received.r, &compareColors.received.g,
					&compareColors.received.b, &compareColors.received.a,
				)
				shiftPixel(
					&compareColors.expected.r, &compareColors.expected.g,
					&compareColors.expected.b, &compareColors.expected.a,
				)
				if compareColors.expected != compareColors.received {
					t.Errorf(
						"Round-trip mismatch. Got color %d: %+v, expected %+v",
						index, got[index], expectedColor,
					)
				}
			}
		})
	}
}
