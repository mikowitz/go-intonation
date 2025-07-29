package intonation

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestNewRatio(t *testing.T) {
	t.Run("a simple ratio", func(t *testing.T) {
		r := NewRatio(5, 4)

		if r.String() != "5/4" {
			t.Errorf("Expected '5/4', got '%s'", r.String())
		}
	})

	t.Run("normalizing a ratio < 2", func(t *testing.T) {
		r := NewRatio(12, 5)

		if r.String() != "6/5" {
			t.Errorf("Expected '6/5', got '%s'", r.String())
		}
	})

	t.Run("normalizing a ratio > 1", func(t *testing.T) {
		r := NewRatio(3, 4)

		if r.String() != "3/2" {
			t.Errorf("Expected '3/2', got '%s'", r.String())
		}
	})
}

const EPSILON = 0.0001

func TestApproximate12EDOInterval(t *testing.T) {
	type testCase struct {
		description      string
		ratio            Ratio
		expectedInterval Interval
		expectedCents    float64
	}
	testCases := []testCase{
		{"unison", NewRatio(1, 1), Unison, 0.0},
		{"perfect fourth", NewRatio(4, 3), PerfectFourth, -1.955},
		{"augmented fourth", NewRatio(11, 8), AugmentedFourth, -48.6821},
		{"perfect fifth", NewRatio(3, 2), PerfectFifth, 1.955},
		{"minor seventh", NewRatio(7, 4), MinorSeventh, -31.1741},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			a := tt.ratio.Approximate12EDOInterval()

			if a.Interval() != tt.expectedInterval {
				t.Errorf("Expected %q, got %q", tt.expectedInterval, a.Interval())
			}

			if math.Abs(a.CentsOffset()-tt.expectedCents) > EPSILON {
				t.Errorf("Expected %.4f cents offset, got %.4f", tt.expectedCents, a.CentsOffset())
			}
		})
	}
}

func TestApproximateEDOInterval(t *testing.T) {
	type testCase struct {
		description      string
		ratio            Ratio
		edo              uint
		expectedInterval Interval
		expectedCents    float64
	}
	testCases := []testCase{
		{"unison", NewRatio(1, 1), 19, Interval{0, 19}, 0.0},
		{"6/19", NewRatio(5, 4), 19, Interval{6, 19}, 7.3663},
		{"10/13", NewRatio(7, 4), 13, Interval{10, 13}, 45.7490},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			a := tt.ratio.ApproximateEDOInterval(tt.edo)

			if a.Interval() != tt.expectedInterval {
				t.Errorf("Expected %q, got %q", tt.expectedInterval, a.Interval())
			}

			if math.Abs(a.CentsOffset()-tt.expectedCents) > EPSILON {
				t.Errorf("Expected %.4f cents offset, got %.4f", tt.expectedCents, a.CentsOffset())
			}
		})
	}
}

func TestRatioMul(t *testing.T) {
	testCases := []struct {
		lhs, rhs, expected Ratio
	}{
		{NewRatio(1, 1), NewRatio(3, 2), NewRatio(3, 2)},
		{NewRatio(3, 2), NewRatio(3, 2), NewRatio(9, 4)},
		{NewRatio(7, 4), NewRatio(11, 7), NewRatio(11, 8)},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s * %s = %s", tc.lhs, tc.rhs, tc.expected), func(t *testing.T) {
			actual := tc.lhs.Mul(tc.rhs)

			if tc.expected != actual {
				t.Errorf("expected %s * %s = %s, got %s", tc.lhs, tc.rhs, tc.expected, actual)
			}
		})
	}
}

func TestRatioPow(t *testing.T) {
	testCases := []struct {
		base     Ratio
		power    int
		expected Ratio
	}{
		{NewRatio(1, 1), 0, NewRatio(1, 1)},
		{NewRatio(1, 1), 4, NewRatio(1, 1)},
		{NewRatio(3, 2), 0, NewRatio(1, 1)},
		{NewRatio(3, 2), 2, NewRatio(9, 4)},
		{NewRatio(7, 4), 5, NewRatio(16807, 16384)},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s ^ %d = %s", tc.base, tc.power, tc.expected), func(t *testing.T) {
			actual := tc.base.Pow(tc.power)

			if tc.expected != actual {
				t.Errorf("expected %s ^ %d = %s, got %s", tc.base, tc.power, tc.expected, actual)
			}
		})
	}
}

func TestRatioFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected Ratio
		error    bool
	}{
		{"3/2", NewRatio(3, 2), false},
		{"20/15", NewRatio(4, 3), false},
		{"three/2", Ratio{}, true},
		{"three", Ratio{}, true},
		{"4/three", Ratio{}, true},
		{"4/-2", Ratio{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual, err := NewRatioFromString(tc.input)

			if tc.error {
				expectedErr := RatioParseError{tc.input}
				if !reflect.DeepEqual(&expectedErr, err) {
					t.Errorf("expected error %s, got %s", &expectedErr, err)
				}
				return
			}

			if tc.expected != actual {
				t.Errorf("expected ratio %s, got %s", tc.expected, actual)
			}
		})
	}
}
