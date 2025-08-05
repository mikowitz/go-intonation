package internal

import (
	"context"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type BeepAudioOutput struct {
	SampleRate beep.SampleRate
}

func (output BeepAudioOutput) PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error {
	speaker.Init(output.SampleRate, 4800)

	chordTones := []beep.Streamer{}
	for _, f := range frequencies {
		tone, err := generators.SineTone(output.SampleRate, f)
		if err != nil {
			return err
		}
		chordTones = append(chordTones, tone)
	}

	chord := beep.Take(output.SampleRate.N(duration), beep.Mix(chordTones...))
	chord = &effects.Volume{
		Streamer: chord,
		Base:     2,
		Volume:   -3,
	}

	return output.playWithContext(ctx, chord)
}

func (output BeepAudioOutput) PlayTone(ctx context.Context, frequency float64, duration time.Duration) error {
	speaker.Init(output.SampleRate, 4800)

	tone, err := generators.SineTone(output.SampleRate, frequency)
	if err != nil {
		return err
	}

	tone = beep.Take(output.SampleRate.N(duration), tone)
	tone = &effects.Volume{
		Streamer: tone,
		Base:     2,
		Volume:   -3,
	}

	return output.playWithContext(ctx, tone)
}

func (output BeepAudioOutput) playWithContext(ctx context.Context, streamer beep.Streamer) error {
	speaker.Init(output.SampleRate, 4800)
	done := make(chan bool, 1)

	speaker.Play(beep.Seq(
		streamer,
		beep.Callback(func() {
			done <- true
		}),
	))

	select {
	case <-ctx.Done():
		speaker.Clear()
		return ctx.Err()
	case <-done:
		return nil
	}
}
