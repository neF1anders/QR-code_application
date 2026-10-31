package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/neF1anders/qr/internal/domain"
)

func TestLoadWithoutFileReturnsDefaults(t *testing.T) {
	settings, err := NewJSONFile(filepath.Join(t.TempDir(), "config.json")).Load()
	if err != nil || settings != domain.DefaultSettings() {
		t.Fatalf("Load() = %+v, %v; want defaults, nil", settings, err)
	}
}

func TestSaveThenLoadRoundTrip(t *testing.T) {
	file := NewJSONFile(filepath.Join(t.TempDir(), "nested", "config.json"))
	want := domain.Settings{TTL: 5 * time.Minute, SizePx: 512, Recovery: domain.RecoveryHigh}

	if err := file.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := file.Load()
	if err != nil || got != want {
		t.Fatalf("Load() = %+v, %v; want %+v", got, err, want)
	}
}

func TestLoadFillsMissingFieldsWithDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"ttl_seconds": 60}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := NewJSONFile(path).Load()
	want := domain.DefaultSettings()
	want.TTL = time.Minute
	if err != nil || got != want {
		t.Fatalf("Load() = %+v, %v; want %+v", got, err, want)
	}
}

func TestLoadReportsCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewJSONFile(path).Load(); err == nil {
		t.Fatal("Load() must fail on corrupt JSON")
	}
}
