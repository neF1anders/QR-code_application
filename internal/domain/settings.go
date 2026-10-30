// Package domain holds the core rules of the qr utility: settings and QR image lifetime.
package domain

import (
	"errors"
	"fmt"
	"time"
)

// RecoveryLevel is the QR error correction level: higher levels survive more damage but hold less data.
type RecoveryLevel string

const (
	RecoveryLow     RecoveryLevel = "low"
	RecoveryMedium  RecoveryLevel = "medium"
	RecoveryHigh    RecoveryLevel = "high"
	RecoveryHighest RecoveryLevel = "highest"
)

// RecoveryLevels lists every supported level from weakest to strongest.
var RecoveryLevels = []RecoveryLevel{RecoveryLow, RecoveryMedium, RecoveryHigh, RecoveryHighest}

const (
	MinSizePx = 128
	MaxSizePx = 2048
)

var (
	ErrNegativeTTL          = errors.New("ttl must not be negative")
	ErrSizeOutOfRange       = fmt.Errorf("size must be between %d and %d px", MinSizePx, MaxSizePx)
	ErrUnknownRecoveryLevel = errors.New("unknown recovery level")
)

// Settings is the user configuration persisted between sessions.
// A zero TTL means the image is deleted as soon as its window is closed.
type Settings struct {
	TTL      time.Duration
	SizePx   int
	Recovery RecoveryLevel
}

// DefaultSettings returns the configuration used when the user has not changed anything.
func DefaultSettings() Settings {
	return Settings{TTL: 0, SizePx: 384, Recovery: RecoveryMedium}
}

// Validate reports the first rule the settings break, or nil.
func (s Settings) Validate() error {
	if s.TTL < 0 {
		return ErrNegativeTTL
	}
	if s.SizePx < MinSizePx || s.SizePx > MaxSizePx {
		return ErrSizeOutOfRange
	}
	if !s.Recovery.IsKnown() {
		return fmt.Errorf("%w: %q", ErrUnknownRecoveryLevel, s.Recovery)
	}
	return nil
}

// IsKnown reports whether the level is one of RecoveryLevels.
func (l RecoveryLevel) IsKnown() bool {
	for _, known := range RecoveryLevels {
		if l == known {
			return true
		}
	}
	return false
}
