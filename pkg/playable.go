package intonation

import "time"

type AudioOutput interface {
	PlayTone(frequency float64, duration time.Duration) error
	PlayChord(frequencies []float64, duration time.Duration) error
}

type Playable interface {
	Play(output AudioOutput)
	PlayInterval(output AudioOutput)
	PlayChord(output AudioOutput)
}
