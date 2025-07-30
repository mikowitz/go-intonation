package internal

import (
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type BeepAudioOutput struct {
	SampleRate beep.SampleRate
}

func (output BeepAudioOutput) PlayChord(frequencies []float64, duration time.Duration) error {
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
		Base:     2,
		Volume:   -2,
	}

	ch := make(chan struct{})
	speaker.Play(beep.Seq(
		beep.Take(output.SampleRate.N(duration), chord),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	))
	<-ch

	return nil
}

func (output BeepAudioOutput) PlayTone(frequency float64, duration time.Duration) error {
	speaker.Init(output.SampleRate, 4800)

	tone, err := generators.SineTone(output.SampleRate, frequency)
	if err != nil {
		return err
	}

	tone = &effects.Volume{
		Streamer: tone,
		Base:     2,
		Volume:   -3,
	}

	ch := make(chan struct{})
	streamers := []beep.Streamer{
		beep.Take(output.SampleRate.N(duration), tone),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	}
	speaker.Play(beep.Seq(streamers...))
	<-ch

	return nil
}
