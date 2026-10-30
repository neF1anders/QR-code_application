// Package config persists settings as a JSON file.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

type fileFormat struct {
	TTLSeconds    int64  `json:"ttl_seconds"`
	SizePx        int    `json:"size_px"`
	RecoveryLevel string `json:"recovery_level"`
}

// JSONFile stores settings in a single JSON file.
type JSONFile struct {
	path string
}

func NewJSONFile(path string) *JSONFile {
	return &JSONFile{path: path}
}

// Load returns the saved settings, or the defaults when the file does not exist yet.
// Fields missing from the file keep their default values.
func (f *JSONFile) Load() (domain.Settings, error) {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return domain.DefaultSettings(), nil
	}
	if err != nil {
		return domain.Settings{}, err
	}
	stored := toFileFormat(domain.DefaultSettings())
	if err := json.Unmarshal(data, &stored); err != nil {
		return domain.Settings{}, err
	}
	return domain.Settings{
		TTL:      time.Duration(stored.TTLSeconds) * time.Second,
		SizePx:   stored.SizePx,
		Recovery: domain.RecoveryLevel(stored.RecoveryLevel),
	}, nil
}

// Save writes settings atomically so a crash never leaves a half-written config.
func (f *JSONFile) Save(settings domain.Settings) error {
	data, err := json.MarshalIndent(toFileFormat(settings), "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(f.path), ".config-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), f.path)
}

func toFileFormat(settings domain.Settings) fileFormat {
	return fileFormat{
		TTLSeconds:    int64(settings.TTL / time.Second),
		SizePx:        settings.SizePx,
		RecoveryLevel: string(settings.Recovery),
	}
}
