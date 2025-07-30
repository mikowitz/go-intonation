package intonation

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type RatioParseError struct {
	input string
}

func (e *RatioParseError) Error() string {
	return fmt.Sprintf("could not parse ratio %s", e.input)
}

type Ratio struct {
	numer, denom uint
}

func NewRatio(n, d uint) Ratio {
	n, d = normalize(n, d)
	g := gcd(n, d)
	return Ratio{n / g, d / g}
}

func NewRatioFromString(input string) (Ratio, error) {
	parts := strings.Split(input, "/")

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return Ratio{}, &RatioParseError{input}
	}
	d, err := strconv.Atoi(parts[1])
	if err != nil {
		return Ratio{}, &RatioParseError{input}
	}

	if n < 0 || d <= 0 {
		return Ratio{}, &RatioParseError{input}
	}

	return NewRatio(uint(n), uint(d)), nil
}

func (r Ratio) String() string {
	return fmt.Sprintf("%d/%d", r.numer, r.denom)
}

func (r Ratio) Float() float64 {
	return float64(r.numer) / float64(r.denom)
}

func (r Ratio) Mul(rhs Ratio) Ratio {
	return NewRatio(r.numer*rhs.numer, r.denom*rhs.denom)
}

func (r Ratio) Pow(base int) Ratio {
	var n uint = 1
	var d uint = 1

	for range base {
		n *= r.numer
		d *= r.denom
	}

	return NewRatio(n, d)
}

func (r Ratio) Approximate12EDOInterval() ApproximateEDOInterval {
	return r.ApproximateEDOInterval(12)
}

func (r Ratio) ApproximateEDOInterval(edo uint) ApproximateEDOInterval {
	f := r.Float()

	edoStepCents := 1200.0 / float64(edo)

	jiCents := 1200.0 * math.Log2(f)

	etCents := math.Round(jiCents/edoStepCents) * edoStepCents

	return ApproximateEDOInterval{
		Interval{uint(etCents/edoStepCents) % edo, edo},
		jiCents - etCents,
	}
}

func (r Ratio) dyad() []float64 {
	return []float64{MiddleC, MiddleC * r.Float()}
}

func (r Ratio) PlayInterval(output AudioOutput) error {
	dyad := r.dyad()
	err := output.PlayTone(dyad[0], 2*time.Second)
	if err != nil {
		return err
	}
	return output.PlayTone(dyad[1], 2*time.Second)
}

func (r Ratio) PlayChord(output AudioOutput) error {
	return output.PlayChord(r.dyad(), 2*time.Second)
}

func (r Ratio) Play(output AudioOutput) error {
	err := r.PlayInterval(output)
	if err != nil {
		return err
	}
	return r.PlayChord(output)
}

func gcd(a, b uint) uint {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func normalize(a, b uint) (uint, uint) {
	f := a / b
	if f >= 2 {
		return normalize(a, b*2)
	}
	if f < 1 {
		return normalize(a*2, b)
	}
	return a, b
}
