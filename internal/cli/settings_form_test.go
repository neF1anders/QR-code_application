package cli

import (
	"testing"
	"time"
)

func TestTTLChoiceFor(t *testing.T) {
	tests := []struct {
		ttl  time.Duration
		want time.Duration
	}{
		{0, 0},
		{time.Hour, time.Hour},
		{42 * time.Second, customTTL},
	}
	for _, tt := range tests {
		if got := ttlChoiceFor(tt.ttl); got != tt.want {
			t.Errorf("ttlChoiceFor(%v) = %v, want %v", tt.ttl, got, tt.want)
		}
	}
}

func TestParseSeconds(t *testing.T) {
	if got, err := parseSeconds("90"); err != nil || got != 90*time.Second {
		t.Fatalf("parseSeconds(90) = %v, %v", got, err)
	}
	for _, input := range []string{"", "-1", "1.5", "abc", "99999999999999999"} {
		if _, err := parseSeconds(input); err == nil {
			t.Errorf("parseSeconds(%q) must fail", input)
		}
	}
}

func TestSizeOptionsKeepsCustomCurrentSize(t *testing.T) {
	options := sizeOptions(1000)
	if len(options) != len(sizePresets)+1 || options[0].Value != 1000 {
		t.Fatalf("sizeOptions(1000) = %v", options)
	}
	if len(sizeOptions(sizePresets[0])) != len(sizePresets) {
		t.Fatal("a preset size must not be duplicated")
	}
}
