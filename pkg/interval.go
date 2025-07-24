package intonation

import "math"

type Interval struct {
	steps, edo uint
}

func NewInterval(steps, edo uint) Interval {
	return Interval{steps, edo}
}

func (i Interval) Cents() float64 {
	return float64(i.steps) * 1200.0 / float64(i.edo)
}

func (i Interval) Play() {
	stepRatio := math.Pow(2, 1.0/float64(i.edo))
	intervalRatio := math.Pow(stepRatio, float64(i.steps))

	dyad{256.0, 256.0 * intervalRatio}.Play()
}
