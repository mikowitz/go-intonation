package intonation

import "errors"

var (
	ErrInvalidRatio      = errors.New("invalid ratio format")
	ErrLatticeDimensions = errors.New("too many access indices passed")
)
