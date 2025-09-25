package audio

import (
	"context"
	"time"

	"github.com/gopxl/beep/v2"
)

type AudioOutput interface {
	PlayTone(ctx context.Context, frequency float64, duration time.Duration) error
	PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error
	StreamChord(ctx context.Context, frequencies []float64) (*beep.Streamer, error)
	PauseStream(ctx context.Context, stream *beep.Streamer)
}
