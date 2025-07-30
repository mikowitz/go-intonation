package intonation

import "time"

const MiddleC float64 = 256.0

type AudioOutput interface {
	PlayTone(frequency float64, duration time.Duration) error
	PlayChord(frequencies []float64, duration time.Duration) error
}

type Playable interface {
	Play(output AudioOutput) error
	PlayInterval(output AudioOutput) error
	PlayChord(output AudioOutput) error
}
