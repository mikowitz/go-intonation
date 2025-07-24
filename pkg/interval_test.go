package intonation

import "testing"

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
