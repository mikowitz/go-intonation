package intonation

import "testing"

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
