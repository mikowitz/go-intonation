package audio

import (
	"context"
	"time"
)

type AudioOutput interface {
	PlayTone(ctx context.Context, frequency float64, duration time.Duration) error
	PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error
}
