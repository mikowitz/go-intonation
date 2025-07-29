package intonation

import "testing"

func TestNewLattice(t *testing.T) {
	l := NewLattice(
		NewRatio(3, 2),
		NewRatio(5, 4),
		NewRatio(7, 4),
	)

	if len(l) != 3 {
		t.Errorf("Expected a lattice with 3 dimensions, got %d", len(l))
	}
}

func TestLatticeAt(t *testing.T) {
	testCases := []struct {
		ratios   []Ratio
		access   []int
		expected Ratio
	}{
		{[]Ratio{NewRatio(3, 2), NewRatio(5, 4)}, []int{0, 0}, NewRatio(1, 1)},
		{[]Ratio{NewRatio(3, 2), NewRatio(5, 4), NewRatio(8, 7)}, []int{1, 2, 3}, NewRatio(600, 343)},
	}

	for _, tc := range testCases {
		t.Run(tc.expected.String(), func(t *testing.T) {
			l := NewLattice(tc.ratios...)
			actual, err := l.At(tc.access...)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			if tc.expected != actual {
				t.Errorf("Expected %s, got %s", tc.expected, actual)
			}
		})
	}

	t.Run("errors if too many dimensions are passed", func(t *testing.T) {
		l := NewLattice(NewRatio(3, 2), NewRatio(5, 4))

		_, err := l.At(1, 2, 3)

		if err != ErrLatticeDimensions {
			t.Errorf("expected ErrLatticeDimensions, got %q", err)
		}
	})

	t.Run("sets unset dimension access to 0", func(t *testing.T) {
		l := NewLattice(NewRatio(3, 2), NewRatio(5, 4), NewRatio(11, 4))

		actual, err := l.At(1, 2)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		expected := NewRatio(75, 64)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})
}
