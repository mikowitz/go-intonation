package intonation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var ErrInvalidRatio = errors.New("invalid ratio format")

type Ratio struct {
	Numer, Denom uint
}

func NewRatio(n, d uint) Ratio {
	n, d = normalize(n, d)
	g := gcd(n, d)
	return Ratio{n / g, d / g}
}

func NewRatioFromString(input string) (Ratio, error) {
	parts := strings.Split(input, "/")

	if len(parts) != 2 {
		return Ratio{}, fmt.Errorf("%s %q", ErrInvalidRatio, input)
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return Ratio{}, fmt.Errorf("%s %q", ErrInvalidRatio, input)
	}
	d, err := strconv.Atoi(parts[1])
	if err != nil {
		return Ratio{}, fmt.Errorf("%s %q", ErrInvalidRatio, input)
	}

	if n < 0 || d <= 0 {
		return Ratio{}, fmt.Errorf("%s %q", ErrInvalidRatio, input)
	}

	return NewRatio(uint(n), uint(d)), nil
}

func (r Ratio) String() string {
	return fmt.Sprintf("%d/%d", r.Numer, r.Denom)
}

func (r Ratio) Float() float64 {
	return float64(r.Numer) / float64(r.Denom)
}

func (r Ratio) Mul(rhs Ratio) Ratio {
	return NewRatio(r.Numer*rhs.Numer, r.Denom*rhs.Denom)
}

func (r Ratio) Pow(base int) Ratio {
	var n uint = 1
	var d uint = 1

	for range base {
		n *= r.Numer
		d *= r.Denom
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

func (r Ratio) Dyad() []float64 {
	return []float64{MiddleCFrequency, MiddleCFrequency * r.Float()}
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
