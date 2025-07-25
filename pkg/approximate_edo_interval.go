package intonation

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
