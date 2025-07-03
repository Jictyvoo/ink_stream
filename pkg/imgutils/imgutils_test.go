package imgutils

import (
	"reflect"
	"testing"
)

func TestNormalizePixel(t *testing.T) {
	type testCase[T Number] struct {
		name     string
		input    T
		expected uint8
	}

	intTests := []testCase[int]{
		{"Within range", 128, 128},
		{"Below range", -10, 0},
		{"Above range", 300, 255},
		{"At min", 0, 0},
		{"At max", 255, 255},
	}

	floatTests := []testCase[float64]{
		{"Float within range", 128.9, 128},
		{"Float below range", -5.5, 0},
		{"Float above range", 300.1, 255},
		{"Float at min", 0.0, 0},
		{"Float at max", 255.0, 255},
	}

	for _, tt := range intTests {
		t.Run("int/"+tt.name, func(t *testing.T) {
			result := NormalizePixel(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePixel(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}

	for _, tt := range floatTests {
		t.Run("float/"+tt.name, func(t *testing.T) {
			result := NormalizePixel(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePixel(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMargins_UpdateNonEmpty(t *testing.T) {
	type args[T comparable] struct {
		initial Margins[T]
		other   Margins[T]
		want    Margins[T]
	}

	tests := []struct {
		name string
		args args[int]
	}{
		{
			name: "All zero fields get updated",
			args: args[int]{
				initial: Margins[int]{},
				other:   Margins[int]{Top: 1, Bottom: 2, Left: 3, Right: 4},
				want:    Margins[int]{Top: 1, Bottom: 2, Left: 3, Right: 4},
			},
		},
		{
			name: "Some fields are already set, only unset ones are updated",
			args: args[int]{
				initial: Margins[int]{Left: 10, Bottom: 20},
				other:   Margins[int]{Top: 1, Bottom: 2, Left: 3, Right: 4},
				want:    Margins[int]{Top: 1, Bottom: 20, Left: 10, Right: 4},
			},
		},
		{
			name: "No fields updated if all are already set",
			args: args[int]{
				initial: Margins[int]{Top: 5, Bottom: 6, Left: 7, Right: 8},
				other:   Margins[int]{Top: 1, Bottom: 2, Left: 3, Right: 4},
				want:    Margins[int]{Top: 5, Bottom: 6, Left: 7, Right: 8},
			},
		},
		{
			name: "Other has zero values — should not overwrite",
			args: args[int]{
				initial: Margins[int]{Top: 9},
				other:   Margins[int]{},
				want:    Margins[int]{Top: 9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.initial
			got.UpdateNonEmpty(tt.args.other)
			if !reflect.DeepEqual(got, tt.args.want) {
				t.Errorf("UpdateNonEmpty() = %+v; want %+v", got, tt.args.want)
			}
		})
	}
}
