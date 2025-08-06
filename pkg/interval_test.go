package intonation

import (
	"context"
	"errors"
	"math"
	"reflect"
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
		edo              EDO
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
		expectedSteps Steps
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

func TestPrinting12EDOIntervals(t *testing.T) {
	testCases := []struct {
		interval TwelveEDOInterval
		name     string
	}{
		{interval: Unison, name: "Unison"},
		{interval: MinorSecond, name: "Minor Second"},
		{interval: MajorSecond, name: "Major Second"},
		{interval: MinorThird, name: "Minor Third"},
		{interval: MajorThird, name: "Major Third"},
		{interval: PerfectFourth, name: "Perfect Fourth"},
		{interval: AugmentedFourth, name: "Augmented Fourth"},
		{interval: PerfectFifth, name: "Perfect Fifth"},
		{interval: MinorSixth, name: "Minor Sixth"},
		{interval: MajorSixth, name: "Major Sixth"},
		{interval: MinorSeventh, name: "Minor Seventh"},
		{interval: MajorSeventh, name: "Major Seventh"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.interval.String()
			if actual != tc.name {
				t.Errorf("Expected %s, got %s", tc.name, actual)
			}
		})
	}
}

func TestPrintingIntervals(t *testing.T) {
	testCases := []struct {
		interval Interval
		name     string
	}{
		{NewInterval(7, 13), "13-EDO 7 steps"},
		{NewInterval(1, 17), "17-EDO 1 step"},
		{NewInterval(8, 20), "20-EDO 8 steps"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.interval.String()
			if actual != tc.name {
				t.Errorf("Expected %s, got %s", tc.name, actual)
			}
		})
	}
}

func TestIntervalPlay(t *testing.T) {
	i := PerfectFourth
	output := &TestAudioOutput{}

	Play(i, context.Background(), output)

	if len(output.output) != 4 {
		t.Errorf("expected 4 tones played, got %d", len(output.output))
	}
	expected := []TestAudioOutputRecord{
		{f: 256.0, chord: false},
		{f: 341.7190026675288, chord: false},
		{f: 256.0, chord: true},
		{f: 341.7190026675288, chord: true},
	}
	if !reflect.DeepEqual(expected, output.output) {
		t.Errorf("expected\n%v\ngot\n%v", expected, output.output)
	}
}

func TestIntervalPlayError(t *testing.T) {
	i := MajorSecond
	output := IntervalErroringOutput{}
	expected := errors.New("couldn't play tone")
	t.Run("play interval", func(t *testing.T) {
		err := PlayInterval(i, context.Background(), output)

		if expected.Error() != err.Error() {
			t.Errorf("expected %s, got %s", expected, err)
		}
	})

	t.Run("play chord", func(t *testing.T) {
		err := PlayChord(i, context.Background(), output)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})

	t.Run("play", func(t *testing.T) {
		err := Play(i, context.Background(), output)

		if expected.Error() != err.Error() {
			t.Errorf("expected %s, got %s", expected, err)
		}
	})
}

func TestIntervalPlayErrorWithChord(t *testing.T) {
	i := MinorSixth
	output := ChordErroringOutput{}
	expected := errors.New("couldn't play chord")
	t.Run("play interval", func(t *testing.T) {
		err := PlayInterval(i, context.Background(), output)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})

	t.Run("play chord", func(t *testing.T) {
		err := PlayChord(i, context.Background(), output)
		if expected.Error() != err.Error() {
			t.Errorf("expected %s, got %s", expected, err)
		}
	})

	t.Run("play", func(t *testing.T) {
		err := Play(i, context.Background(), output)

		if expected.Error() != err.Error() {
			t.Errorf("expected %s, got %s", expected, err)
		}
	})
}
