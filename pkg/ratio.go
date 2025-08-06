package intonation

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type (
	Numer = uint
	Denom = uint
)

type Ratio struct {
	Numer
	Denom
}

func NewRatio(n, d uint) Ratio {
	n, d = normalize(n, d)
	g := gcd(n, d)
	return Ratio{Numer(n / g), Denom(d / g)}
}

func NewRatioFromString(input string) (Ratio, error) {
	parts := strings.Split(input, "/")

	error := fmt.Errorf("parsing ratio %q: %w", input, ErrInvalidRatio)

	if len(parts) != 2 {
		return Ratio{}, error
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return Ratio{}, error
	}
	d, err := strconv.Atoi(parts[1])
	if err != nil {
		return Ratio{}, error
	}

	if n < 0 || d <= 0 {
		return Ratio{}, error
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
	var n Numer = 1
	var d Denom = 1

	for range base {
		n *= r.Numer
		d *= r.Denom
	}

	return NewRatio(n, d)
}

func (r Ratio) Approximate12EDOInterval() ApproximateEDOInterval {
	return r.ApproximateEDOInterval(12)
}

func (r Ratio) ApproximateEDOInterval(edo EDO) ApproximateEDOInterval {
	f := r.Float()

	edoStepCents := 1200.0 / float64(edo)

	jiCents := 1200.0 * math.Log2(f)

	etCents := math.Round(jiCents/edoStepCents) * edoStepCents

	return ApproximateEDOInterval{
		Interval{Steps(uint(etCents/edoStepCents) % uint(edo)), edo},
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
