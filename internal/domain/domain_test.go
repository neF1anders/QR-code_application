package domain

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultSettingsAreValid(t *testing.T) {
	if err := DefaultSettings().Validate(); err != nil {
		t.Fatalf("default settings invalid: %v", err)
	}
}

func TestSettingsValidate(t *testing.T) {
	valid := DefaultSettings()
	tests := []struct {
		name   string
		modify func(*Settings)
		want   error
	}{
		{"positive ttl", func(s *Settings) { s.TTL = time.Minute }, nil},
		{"negative ttl", func(s *Settings) { s.TTL = -time.Second }, ErrNegativeTTL},
		{"min size", func(s *Settings) { s.SizePx = MinSizePx }, nil},
		{"too small", func(s *Settings) { s.SizePx = MinSizePx - 1 }, ErrSizeOutOfRange},
		{"too large", func(s *Settings) { s.SizePx = MaxSizePx + 1 }, ErrSizeOutOfRange},
		{"unknown level", func(s *Settings) { s.Recovery = "ultra" }, ErrUnknownRecoveryLevel},
		{"empty level", func(s *Settings) { s.Recovery = "" }, ErrUnknownRecoveryLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := valid
			tt.modify(&s)
			if err := s.Validate(); !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestStoredQRIsExpired(t *testing.T) {
	created := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	qr := StoredQR{Path: "x.png", CreatedAt: created}
	tests := []struct {
		name string
		now  time.Time
		ttl  time.Duration
		want bool
	}{
		{"zero ttl is expired immediately", created, 0, true},
		{"within ttl", created.Add(59 * time.Second), time.Minute, false},
		{"exactly at ttl", created.Add(time.Minute), time.Minute, true},
		{"after ttl", created.Add(time.Hour), time.Minute, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := qr.IsExpired(tt.now, tt.ttl); got != tt.want {
				t.Fatalf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
