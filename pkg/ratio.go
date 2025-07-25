package intonation

import (
	"fmt"

	"github.com/mikowitz/intonation/internal"
)

type Ratio struct {
	numer, denom uint
}

func NewRatio(n, d uint) Ratio {
	n, d = normalize(n, d)
	g := gcd(n, d)
	return Ratio{n / g, d / g}
}

func (r Ratio) String() string {
	return fmt.Sprintf("%d/%d", r.numer, r.denom)
}

func (r Ratio) Float() float64 {
	return float64(r.numer) / float64(r.denom)
}

func (r Ratio) Play() {
	internal.Dyad{256.0, 256.0 * r.Float()}.Play()
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
