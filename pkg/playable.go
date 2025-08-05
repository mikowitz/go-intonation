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
	Dyad() []float64
}

func PlayInterval(p Playable, ctx context.Context, output AudioOutput) error {
	dyad := p.Dyad()
	err := output.PlayTone(ctx, dyad[0], 2*time.Second)
	if err != nil {
		return err
	}
	return output.PlayTone(ctx, dyad[1], 2*time.Second)
}

func PlayChord(p Playable, ctx context.Context, output AudioOutput) error {
	return output.PlayChord(ctx, p.Dyad(), 2*time.Second)
}

func Play(p Playable, ctx context.Context, output AudioOutput) error {
	err := PlayInterval(p, ctx, output)
	if err != nil {
		return err
	}
	return PlayChord(p, ctx, output)
}
