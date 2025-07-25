package intonation

import "testing"

func TestApproximateIntervalString(t *testing.T) {
	testCases := []struct {
		approximateInterval ApproximateEDOInterval
		name                string
	}{
		{ApproximateEDOInterval{Unison, 7.234567}, "Unison (+7.2346)"},
		{ApproximateEDOInterval{Unison, 0}, "Unison (+0.0000)"},
		{ApproximateEDOInterval{Unison, 0.00009}, "Unison (+0.0001)"},
		{ApproximateEDOInterval{Unison, -17.10007}, "Unison (-17.1001)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.approximateInterval.String() != tc.name {
				t.Errorf("Expected %s, got %s", tc.name, tc.approximateInterval.String())
			}
		})
	}
}
