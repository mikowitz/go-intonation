package intonation

import (
	"context"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/mikowitz/intonation/pkg/audio"
)

type Playable interface {
	Dyad() []float64
	String() string
}

func PlayInterval(p Playable, ctx context.Context, output audio.AudioOutput) error {
	dyad := p.Dyad()
	err := output.PlayTone(ctx, dyad[0], 2*time.Second)
	if err != nil {
		return err
	}
	return output.PlayTone(ctx, dyad[1], 2*time.Second)
}

func PlayChord(p Playable, ctx context.Context, output audio.AudioOutput) error {
	return output.PlayChord(ctx, p.Dyad(), 2*time.Second)
}

func Play(p Playable, ctx context.Context, output audio.AudioOutput) error {
	err := PlayInterval(p, ctx, output)
	if err != nil {
		return err
	}
	return PlayChord(p, ctx, output)
}

func StreamChord(p Playable, ctx context.Context, output audio.AudioOutput) (*beep.Ctrl, error) {
	return output.StreamChord(ctx, p.Dyad())
}

func PauseStream(stream *beep.Ctrl, ctx context.Context, output audio.AudioOutput) {
	output.PauseStream(ctx, stream)
}
