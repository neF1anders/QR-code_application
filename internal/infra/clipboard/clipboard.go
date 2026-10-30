// Package clipboard reads text from the system clipboard (Wayland data-control or X11/XWayland).
package clipboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.design/x/clipboard"
)

// readTimeout bounds a read so a stuck clipboard owner cannot hang the CLI.
const readTimeout = 3 * time.Second

// ErrEmpty is returned when the clipboard holds no text.
var ErrEmpty = errors.New("clipboard has no text")

// System is the clipboard of the current desktop session.
type System struct{}

// ReadText returns the clipboard text.
func (System) ReadText() (string, error) {
	if err := clipboard.Init(); err != nil {
		return "", fmt.Errorf("open clipboard: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	data, err := clipboard.Read(ctx, clipboard.FmtText)
	if errors.Is(err, clipboard.ErrNoData) {
		return "", ErrEmpty
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}
