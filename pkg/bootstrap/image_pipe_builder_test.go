package bootstrap

import (
	"testing"

	"github.com/Jictyvoo/ink_stream/internal/imageparser"
	"github.com/Jictyvoo/ink_stream/internal/imageparser/imgpipesteps"
	"github.com/Jictyvoo/ink_stream/pkg/deviceprof"
)

func isTypeOf[T imageparser.PipeStep]() func(inputVal any) bool {
	return func(inputVal any) bool {
		switch inputVal.(type) {
		case T, *T:
			return true
		default:
			return false
		}
	}
}

func TestBuildPipeline(t *testing.T) {
	testCases := []struct {
		name          string
		opts          Options
		expectError   bool
		expectedSteps []func(inputVal any) bool // Types of steps expected in order
	}{
		{
			name: "Successful pipeline with all options",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with rotation enabled",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   true,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with margins disabled",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   false,
				AddMargins:    false,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with grayscale conversion",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  false,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepGrayScaleImage](),
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with rotation and grayscale",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   true,
				AddMargins:    true,
				ColoredPages:  false,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepGrayScaleImage](),
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with CropBasic crop level",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropBasic,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with CropAggressive crop level",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "ltr",
				CropLevel:     CropAggressive,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
		{
			name: "Pipeline with invalid device",
			opts: Options{
				TargetDevice:  "invalid_device",
				ReadDirection: "ltr",
				CropLevel:     CropNormal,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: true,
		},
		{
			name: "Pipeline with empty read direction",
			opts: Options{
				TargetDevice:  deviceprof.DeviceOther,
				ReadDirection: "",
				CropLevel:     CropNormal,
				RotateImage:   false,
				AddMargins:    true,
				ColoredPages:  true,
			},
			expectError: false,
			expectedSteps: []func(inputVal any) bool{
				isTypeOf[*imgpipesteps.StepAutoCropImage](),
				isTypeOf[*imgpipesteps.StepMarginWrapImage](),
				isTypeOf[*imgpipesteps.StepCropOrRotateImage](),
				isTypeOf[*imgpipesteps.StepRescaleImage](),
				isTypeOf[*imgpipesteps.StepAutoContrastImage](),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pipeline, err := BuildPipeline(tc.opts)

			if err != nil {
				if !tc.expectError {
					t.Fatalf("Unexpected error: %v", err.Error())
				}
				return
			} else if tc.expectError {
				t.Fatalf("Unexpected success")
			}

			steps := pipeline.PipeSteps()
			if len(steps) == 0 {
				t.Fatalf("No steps found")
			}

			for index, expectedType := range tc.expectedSteps {
				if !expectedType(steps[index]) {
					t.Errorf("Step at index %d should be of type %T", index, expectedType)
				}
			}
		})
	}
}
