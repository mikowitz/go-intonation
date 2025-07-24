package intonation

import (
	"log"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type Playable interface {
	Play()
}

type IntonationErrorType string

const NewSineToneError IntonationErrorType = "Could not create a new sine tone generator"

type IntonationError struct {
	name IntonationErrorType
	err  error
}

type dyad [2]float64

const (
	sr = beep.SampleRate(48000)
)

var two = sr.N(2 * time.Second)

func (d dyad) Play() {
	speaker.Init(sr, 4800)

	root := lowerVolume(newSineTone(d[0]))
	top := lowerVolume(newSineTone(d[1]))

	chord := effects.Volume{
		Streamer: beep.Mix(root, top),
		Base:     2,
		Volume:   -2,
	}

	ch := make(chan struct{})
	sounds := []beep.Streamer{
		beep.Take(two, root),
		beep.Take(two, top),
		beep.Take(two*3/2, &chord),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	}

	speaker.Play(beep.Seq(sounds...))
	<-ch
}

func newSineTone(freq float64) beep.Streamer {
	s, err := generators.SineTone(sr, freq)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func lowerVolume(s beep.Streamer) beep.Streamer {
	return &effects.Volume{
		Streamer: s,
		Base:     2,
		Volume:   -4,
	}
}
