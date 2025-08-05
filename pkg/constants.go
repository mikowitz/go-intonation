package intonation

import "time"

const (
	// MiddleCFrequency represents the approximate frequency
	// of middle C (C4) in Hz. Uses 256 Hz instead of the more
	// exact 261.63 Hz for mathematical convenience
	MiddleCFrequency float64 = 256.0
	// DefaultPlaybackDuration is the standard duration for
	// audio playback of ratios and frequencies
	DefaultPlaybackDuration time.Duration = 2 * time.Second
)
