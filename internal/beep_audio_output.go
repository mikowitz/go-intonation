package internal

import (
	"context"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type BeepAudioOutput struct {
	SampleRate beep.SampleRate
	once       sync.Once
}

func (output *BeepAudioOutput) initSpeaker() {
	output.once.Do(func() {
		speaker.Init(output.SampleRate, 4800)
	})
}

func (output *BeepAudioOutput) PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error {
	chordTones := []beep.Streamer{}
	for _, f := range frequencies {
		tone, err := generators.SineTone(output.SampleRate, f)
		if err != nil {
			return err
		}
		chordTones = append(chordTones, tone)
	}

	chord := beep.Take(output.SampleRate.N(duration), beep.Mix(chordTones...))

	return output.playWithContext(ctx, chord)
}

func (output *BeepAudioOutput) PlayTone(ctx context.Context, frequency float64, duration time.Duration) error {
	tone, err := generators.SineTone(output.SampleRate, frequency)
	if err != nil {
		return err
	}

	tone = beep.Take(output.SampleRate.N(duration), tone)
	return output.playWithContext(ctx, tone)
}

func applyVolumeEffect(streamer beep.Streamer) beep.Streamer {
	return &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -3,
	}
}

func (output *BeepAudioOutput) playWithContext(ctx context.Context, streamer beep.Streamer) error {
	output.initSpeaker()
	done := make(chan bool, 1)

	speaker.Play(beep.Seq(
		applyVolumeEffect(streamer),
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
