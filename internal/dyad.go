package internal

import (
	"log"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type Dyad [2]float64

func (d Dyad) Play() {
	sr := beep.SampleRate(48000)
	two := sr.N(2 * time.Second)
	speaker.Init(sr, 4800)

	root := lowerVolume(newSineTone(sr, d[0]))
	top := lowerVolume(newSineTone(sr, d[1]))

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

func newSineTone(sr beep.SampleRate, freq float64) beep.Streamer {
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
		Volume:   -3,
	}
}
