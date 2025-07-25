package intonation

import (
	"math"

	"github.com/mikowitz/intonation/internal"
)

type TwelveEDOInterval = Interval

var (
	Unison          = Interval{0, 12}
	MinorSecond     = Interval{1, 12}
	MajorSecond     = Interval{2, 12}
	MinorThird      = Interval{3, 12}
	MajorThird      = Interval{4, 12}
	PerfectFourth   = Interval{5, 12}
	AugmentedFourth = Interval{6, 12}
	PerfectFifth    = Interval{7, 12}
	MinorSixth      = Interval{8, 12}
	MajorSixth      = Interval{9, 12}
	MinorSeventh    = Interval{10, 12}
	MajorSeverth    = Interval{11, 12}
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

func (i Interval) Approximate12EDOInterval() ApproximateEDOInterval {
	return i.ApproximateEDOInterval(12)
}

func (i Interval) ApproximateEDOInterval(edo uint) ApproximateEDOInterval {
	sourceCents := i.Cents()

	targetStepCents := 1200.0 / float64(edo)
	targetCents := math.Round(sourceCents/targetStepCents) * targetStepCents

	return ApproximateEDOInterval{
		Interval{uint(targetCents / targetStepCents), edo},
		sourceCents - targetCents,
	}
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
