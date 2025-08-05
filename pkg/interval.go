package intonation

import (
	"fmt"
	"math"
)

type TwelveEDOInterval = Interval

var (
	Unison          TwelveEDOInterval = Interval{0, 12}
	MinorSecond     TwelveEDOInterval = Interval{1, 12}
	MajorSecond     TwelveEDOInterval = Interval{2, 12}
	MinorThird      TwelveEDOInterval = Interval{3, 12}
	MajorThird      TwelveEDOInterval = Interval{4, 12}
	PerfectFourth   TwelveEDOInterval = Interval{5, 12}
	AugmentedFourth TwelveEDOInterval = Interval{6, 12}
	PerfectFifth    TwelveEDOInterval = Interval{7, 12}
	MinorSixth      TwelveEDOInterval = Interval{8, 12}
	MajorSixth      TwelveEDOInterval = Interval{9, 12}
	MinorSeventh    TwelveEDOInterval = Interval{10, 12}
	MajorSeventh    TwelveEDOInterval = Interval{11, 12}
)

type Interval struct {
	steps, edo uint
}

func NewInterval(steps, edo uint) Interval {
	return Interval{steps, edo}
}

func (i Interval) String() string {
	if i.edo == 12 {
		s, err := i.twelveEDOIntervalString()
		if err != nil {
			return err.Error()
		}
		return s
	}
	noun := "steps"
	if i.steps == 1 {
		noun = "step"
	}
	return fmt.Sprintf("%d-EDO %d %s", i.edo, i.steps, noun)
}

func (i Interval) twelveEDOIntervalString() (string, error) {
	switch i {
	case Unison:
		return "Unison", nil
	case MinorSecond:
		return "Minor Second", nil
	case MajorSecond:
		return "Major Second", nil
	case MinorThird:
		return "Minor Third", nil
	case MajorThird:
		return "Major Third", nil
	case PerfectFourth:
		return "Perfect Fourth", nil
	case AugmentedFourth:
		return "Augmented Fourth", nil
	case PerfectFifth:
		return "Perfect Fifth", nil
	case MinorSixth:
		return "Minor Sixth", nil
	case MajorSixth:
		return "Major Sixth", nil
	case MinorSeventh:
		return "Minor Seventh", nil
	case MajorSeventh:
		return "Major Seventh", nil
	default:
		return "", fmt.Errorf("invalid 12-EDO interval: %v", i)
	}
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

func (i Interval) Dyad() []float64 {
	stepRatio := math.Pow(2, 1.0/float64(i.edo))
	intervalRatio := math.Pow(stepRatio, float64(i.steps))

	return []float64{MiddleC, MiddleC * intervalRatio}
}
