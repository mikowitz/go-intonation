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

const sr = beep.SampleRate(48000)

var two = sr.N(2 * time.Second)

func (d Dyad) PlayInterval() {
	speaker.Init(sr, 4800)

	root := lowerVolume(newSineTone(sr, d[0]))
	top := lowerVolume(newSineTone(sr, d[1]))

	playStreamers(beep.Take(two, root), beep.Take(two, top))
}

func (d Dyad) PlayChord() {
	speaker.Init(sr, 4800)

	root := newSineTone(sr, d[0])
	top := newSineTone(sr, d[1])

	chord := effects.Volume{
		Streamer: beep.Mix(root, top),
		Base:     2,
		Volume:   -2,
	}

	playStreamers(
		beep.Take(two, &chord),
	)
}

func (d Dyad) Play() {
	d.PlayInterval()
	d.PlayChord()
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

func playStreamers(streamers ...beep.Streamer) {
	ch := make(chan struct{})
	streamers = append(streamers, beep.Callback(func() {
		ch <- struct{}{}
	}))
	speaker.Play(beep.Seq(streamers...))
	<-ch
}
