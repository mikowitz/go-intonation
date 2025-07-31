package intonation

import (
	"context"
	"time"
)

const MiddleC float64 = 256.0

type AudioOutput interface {
	PlayTone(ctx context.Context, frequency float64, duration time.Duration) error
	PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error
}

type Playable interface {
	Play(ctx context.Context, output AudioOutput) error
	PlayInterval(ctx context.Context, output AudioOutput) error
	PlayChord(ctx context.Context, output AudioOutput) error
}
