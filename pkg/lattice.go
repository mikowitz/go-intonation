package intonation

import "fmt"

type Lattice []Ratio

func NewLattice(ratios ...Ratio) Lattice {
	l := Lattice{}
	for _, r := range ratios {
		l = append(l, r)
	}
	return l
}

func (l Lattice) At(access ...int) (Ratio, error) {
	if len(access) > len(l) {
		return Ratio{}, fmt.Errorf("accessing lattice: %w", ErrLatticeDimensions)
	}

	r := NewRatio(1, 1)

	for i, a := range access {
		r = r.Mul(l[i].Pow(a))
	}
	return r, nil
}
