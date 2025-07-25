package intonation

import (
	"math"
	"testing"
)

func TestIntervalCents(t *testing.T) {
	t.Run("12 EDO unison", func(t *testing.T) {
		i := NewInterval(0, 12)

		if i.Cents() != 0.0 {
			t.Errorf("Expected cents of 0, got %.3f", i.Cents())
		}
	})

	t.Run("12 EDO major 3rd", func(t *testing.T) {
		i := NewInterval(4, 12)

		if i.Cents() != 400.0 {
			t.Errorf("Expected cents of 400, got %.3f", i.Cents())
		}
	})

	t.Run("19 EDO 7 steps", func(t *testing.T) {
		i := NewInterval(7, 19)

		if (i.Cents() - 442.1053) > 0.00001 {
			t.Errorf("Expected cents of 442.1053, got %.4f", i.Cents())
		}
	})
}

func TestApproximateEDOIntervalFromInterval(t *testing.T) {
	type testCase struct {
		description      string
		interval         Interval
		edo              uint
		expectedInterval Interval
		expectedCents    float64
	}

	testCases := []testCase{
		{"unison 12 <-> 19", NewInterval(0, 12), 19, Interval{0, 19}, 0.0},
		{"unison 15 <-> 19", NewInterval(0, 15), 19, Interval{0, 19}, 0.0},
		{"major 3rd 12 <-> 17", NewInterval(4, 12), 17, Interval{6, 17}, -23.5294},
		{"major 6th 12 <-> 12", NewInterval(9, 12), 12, Interval{9, 12}, 0.0},
		{"6/13 <-> 23", NewInterval(6, 13), 23, Interval{11, 23}, -20.0669},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			a := tt.interval.ApproximateEDOInterval(tt.edo)

			if a.Interval() != tt.expectedInterval {
				t.Errorf("Expected %q, got %q", tt.expectedInterval, a.Interval())
			}

			if math.Abs(a.CentsOffset()-tt.expectedCents) > EPSILON {
				t.Errorf("Expected %.4f cents offset, got %.4f", tt.expectedCents, a.CentsOffset())
			}
		})
	}
}

func TestApproximate12EDOIntervalFromInterval(t *testing.T) {
	type testCase struct {
		description   string
		interval      Interval
		expectedSteps uint
		expectedCents float64
	}

	testCases := []testCase{
		{"6/13 -> 12", NewInterval(6, 13), 6, -46.1538},
		{"6/19 -> 12", NewInterval(6, 19), 4, -21.0526},
		{"6/18 -> 12", NewInterval(6, 18), 4, 0.0},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			expected := NewInterval(tt.expectedSteps, 12)
			a := tt.interval.Approximate12EDOInterval()

			if a.Interval() != expected {
				t.Errorf("Expected %q, got %q", expected, a.Interval())
			}

			if math.Abs(a.CentsOffset()-tt.expectedCents) > EPSILON {
				t.Errorf("Expected %.4f cents offset, got %.4f", tt.expectedCents, a.CentsOffset())
			}
		})
	}
}
