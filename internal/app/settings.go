package app

import (
	"fmt"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

// ManageSettings loads and saves the user settings.
type ManageSettings struct {
	repository SettingsRepository
}

func NewManageSettings(repository SettingsRepository) *ManageSettings {
	return &ManageSettings{repository: repository}
}

// Load returns the saved settings. If they cannot be read or are invalid,
// it returns the defaults together with the reason, so the caller can warn and carry on.
func (u *ManageSettings) Load() (domain.Settings, error) {
	settings, err := u.repository.Load()
	if err != nil {
		return domain.DefaultSettings(), fmt.Errorf("load settings, using defaults: %w", err)
	}
	if err := settings.Validate(); err != nil {
		return domain.DefaultSettings(), fmt.Errorf("invalid settings, using defaults: %w", err)
	}
	return settings, nil
}

// Save validates and persists settings.
func (u *ManageSettings) Save(settings domain.Settings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	if err := u.repository.Save(settings); err != nil {
		return fmt.Errorf("save settings: %w", err)
	}
	return nil
}
