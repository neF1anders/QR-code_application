package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/neF1anders/qr/internal/domain"
)

// ErrNothingToEncode is returned when the text to encode is empty or whitespace only.
var ErrNothingToEncode = errors.New("nothing to encode: text is empty")

// ShowQR encodes text into a QR code, stores it and shows it in a window.
type ShowQR struct {
	clipboard Clipboard
	encoder   Encoder
	store     ImageStore
	viewer    Viewer
	now       Clock
}

func NewShowQR(clipboard Clipboard, encoder Encoder, store ImageStore, viewer Viewer, now Clock) *ShowQR {
	return &ShowQR{clipboard: clipboard, encoder: encoder, store: store, viewer: viewer, now: now}
}

// FromClipboard shows the current clipboard text as a QR code.
func (u *ShowQR) FromClipboard(settings domain.Settings) error {
	text, err := u.clipboard.ReadText()
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}
	return u.FromText(text, settings)
}

// FromText shows text as a QR code. With a zero TTL the image is deleted once the window closes.
func (u *ShowQR) FromText(text string, settings domain.Settings) error {
	if strings.TrimSpace(text) == "" {
		return ErrNothingToEncode
	}
	png, err := u.encoder.Encode(text, settings.SizePx, settings.Recovery)
	if err != nil {
		return fmt.Errorf("encode qr: %w", err)
	}
	stored, err := u.store.Save(png, u.now())
	if err != nil {
		return fmt.Errorf("save qr image: %w", err)
	}

	showErr := u.viewer.Show(png, text)
	if showErr != nil {
		showErr = fmt.Errorf("show qr: %w", showErr)
	}
	var deleteErr error
	if settings.TTL == 0 {
		if err := u.store.Delete(stored); err != nil {
			deleteErr = fmt.Errorf("delete qr image: %w", err)
		}
	}
	return errors.Join(showErr, deleteErr)
}
