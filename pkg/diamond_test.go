package intonation

import (
	"reflect"
	"testing"
)

func TestNewDiamond(t *testing.T) {
	d := NewDiamond(1, 5, 3)

	testCases := []struct {
		row   int
		col   int
		ratio Ratio
	}{
		{0, 0, NewRatio(1, 1)},
		{0, 1, NewRatio(5, 4)},
		{0, 2, NewRatio(3, 2)},
		{1, 0, NewRatio(8, 5)},
		{1, 1, NewRatio(1, 1)},
		{1, 2, NewRatio(6, 5)},
		{2, 0, NewRatio(4, 3)},
		{2, 1, NewRatio(5, 3)},
		{2, 2, NewRatio(1, 1)},
	}

	for _, tc := range testCases {
		if !reflect.DeepEqual(tc.ratio, d[tc.row][tc.col]) {
			t.Errorf("expected %q at (%d, %d), got %q", tc.ratio, tc.row, tc.col, d[tc.row][tc.col])
		}
	}
}

func TestDiamondString(t *testing.T) {
	t.Run("diamond", func(t *testing.T) {
		d := NewDiamond(1, 5, 3)

		expected := `		3/2

	5/4		6/5

1/1		1/1		1/1

	8/5		5/3

		4/3`

		if d.String(formatDiamond) != expected {
			t.Errorf("expected\n\n%q\ngot\n\n%q", expected, d.String(formatDiamond))
		}
	})
	t.Run("square", func(t *testing.T) {
		d := NewDiamond(1, 5, 3)

		expected := `1/1	5/4	3/2

8/5	1/1	6/5

4/3	5/3	1/1`

		if d.String(formatSquare) != expected {
			t.Errorf("expected\n\n%s\ngot\n\n%s", expected, d.String(formatSquare))
		}
	})
}
