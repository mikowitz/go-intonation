package internal

import (
	"context"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

const (
	VolumeBase   = 2
	VolumeVolume = -3
)

type BeepAudioOutput struct {
	SampleRate beep.SampleRate
}

func (output BeepAudioOutput) PlayChord(ctx context.Context, frequencies []float64, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	chordTones := []beep.Streamer{}

	for _, f := range frequencies {
		tone, err := generators.SineTone(output.SampleRate, f)
		if err != nil {
			return err
		}
		chordTones = append(chordTones, tone)
	}

	chord := &effects.Volume{
		Streamer: beep.Mix(chordTones...),
		Base:     VolumeBase,
		Volume:   VolumeVolume,
	}

	ch := make(chan struct{})

	speaker.Play(beep.Seq(
		beep.Take(output.SampleRate.N(duration), chord),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	))

	select {
	case <-ctx.Done():
		speaker.Clear()
		return ctx.Err()
	case <-ch:
		return nil
	}
}

func (output BeepAudioOutput) PlayTone(ctx context.Context, frequency float64, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	speaker.Init(output.SampleRate, 4800)

	tone, err := generators.SineTone(output.SampleRate, frequency)
	if err != nil {
		return err
	}

	tone = &effects.Volume{
		Streamer: tone,
		Base:     VolumeBase,
		Volume:   VolumeVolume,
	}

	ch := make(chan struct{})

	streamers := []beep.Streamer{
		beep.Take(output.SampleRate.N(duration), tone),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	}
	speaker.Play(beep.Seq(streamers...))

	select {
	case <-ctx.Done():
		speaker.Clear()
		return ctx.Err()
	case <-ch:
		return nil
	}
}
