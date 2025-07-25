package intonation

import (
	"math"

	"github.com/mikowitz/intonation/internal"
)

type Interval struct {
	steps, edo uint
}

func NewInterval(steps, edo uint) Interval {
	return Interval{steps, edo}
}

func (i Interval) Cents() float64 {
	return float64(i.steps) * 1200.0 / float64(i.edo)
}

func (i Interval) dyad() internal.Dyad {
	stepRatio := math.Pow(2, 1.0/float64(i.edo))
	intervalRatio := math.Pow(stepRatio, float64(i.steps))

	return internal.Dyad{256.0, 256.0 * intervalRatio}
}

func (i Interval) PlayInterval() {
	i.dyad().PlayInterval()
}

func (i Interval) PlayChord() {
	i.dyad().PlayChord()
}

func (i Interval) Play() {
	i.dyad().Play()
}
