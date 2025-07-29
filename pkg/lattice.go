package intonation

import "errors"

type Lattice []Ratio

var ErrLatticeDimensions = errors.New("too many access indices passed")

func NewLattice(ratios ...Ratio) Lattice {
	l := Lattice{}
	for _, r := range ratios {
		l = append(l, r)
	}
	return l
}

func (l Lattice) At(access ...int) (Ratio, error) {
	if len(access) > len(l) {
		return Ratio{}, ErrLatticeDimensions
	}

	r := NewRatio(1, 1)

	for i, a := range access {
		r = r.Mul(l[i].Pow(a))
	}
	return r, nil
}
