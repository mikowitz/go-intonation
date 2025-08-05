package intonation

import "fmt"

type ApproximateEDOInterval struct {
	interval    Interval
	centsOffset float64
}

func (a ApproximateEDOInterval) Interval() Interval {
	return a.interval
}

func (a ApproximateEDOInterval) CentsOffset() float64 {
	return a.centsOffset
}

func (a ApproximateEDOInterval) String() string {
	sign := "+"
	if a.centsOffset < 0 {
		sign = ""
	}
	return fmt.Sprintf("%s (%s%.4f)", a.interval, sign, a.centsOffset)
}

func (a ApproximateEDOInterval) Dyad() []float64 {
	return a.interval.Dyad()
}
