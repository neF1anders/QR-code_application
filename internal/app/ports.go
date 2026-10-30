// Package app contains the use cases of the qr utility and the ports they depend on.
package app

import (
	"time"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

// Clipboard reads the current text content of the system clipboard.
type Clipboard interface {
	ReadText() (string, error)
}

// Encoder renders text as a PNG QR code.
type Encoder interface {
	Encode(text string, sizePx int, recovery domain.RecoveryLevel) ([]byte, error)
}

// ImageStore keeps rendered QR images on disk.
type ImageStore interface {
	Save(png []byte, createdAt time.Time) (domain.StoredQR, error)
	List() ([]domain.StoredQR, error)
	Delete(qr domain.StoredQR) error
}

// Viewer shows a QR image to the user and blocks until the window is closed.
type Viewer interface {
	Show(png []byte, caption string) error
}

// SettingsRepository persists settings between sessions.
// Load returns domain.DefaultSettings when nothing has been saved yet.
type SettingsRepository interface {
	Load() (domain.Settings, error)
	Save(settings domain.Settings) error
}

// Clock returns the current time; injected so lifetimes are testable.
type Clock func() time.Time
