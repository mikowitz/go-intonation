package internal

import (
	"context"
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
)

func TestBeepAudioOutputContextCancellation(t *testing.T) {
	output := &BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
	
	t.Run("PlayTone cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		
		start := time.Now()
		err := output.PlayTone(ctx, 440.0, 2*time.Second)
		elapsed := time.Since(start)
		
		if err != context.DeadlineExceeded {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
		
		if elapsed > 1*time.Second {
			t.Errorf("expected cancellation within ~500ms, took %v", elapsed)
		}
	})
	
	t.Run("PlayChord cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		
		start := time.Now()
		err := output.PlayChord(ctx, []float64{440.0, 554.37}, 2*time.Second)
		elapsed := time.Since(start)
		
		if err != context.DeadlineExceeded {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
		
		if elapsed > 1*time.Second {
			t.Errorf("expected cancellation within ~500ms, took %v", elapsed)
		}
	})
	
	t.Run("PlayTone with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := output.PlayTone(ctx, 440.0, 2*time.Second)
		
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
	
	t.Run("PlayChord with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := output.PlayChord(ctx, []float64{440.0, 554.37}, 2*time.Second)
		
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}