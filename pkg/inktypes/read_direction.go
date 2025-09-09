package inktypes

import "strings"

type ReadDirection uint8

const (
	ReadUnknown ReadDirection = iota
	ReadLeftToRight
	ReadRightToLeft
)

func NewReadDirection(value string) ReadDirection {
	switch strings.ToLower(value) {
	case ReadRightToLeft.String():
		return ReadRightToLeft
	case ReadLeftToRight.String():
		return ReadLeftToRight
	}

	return ReadUnknown
}

func (rd ReadDirection) String() string {
	switch rd {
	case ReadRightToLeft:
		return "rtl"
	default: // ReadLeftToRight
		return "ltr"
	}
}
