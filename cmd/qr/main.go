// Command qr shows the clipboard text, or the given text, as a QR code in a window.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/neF1anders/qr/internal/app"
	"github.com/neF1anders/qr/internal/cli"
	"github.com/neF1anders/qr/internal/infra/clipboard"
	"github.com/neF1anders/qr/internal/infra/config"
	"github.com/neF1anders/qr/internal/infra/encoder"
	"github.com/neF1anders/qr/internal/infra/storage"
	"github.com/neF1anders/qr/internal/infra/viewer"
)

// version is set at release time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fail(err)
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fail(err)
	}

	store := storage.NewDirectory(filepath.Join(cacheDir, "qr"))
	settingsFile := config.NewJSONFile(filepath.Join(configDir, "qr", "config.json"))

	command := &cli.CLI{
		Settings:     app.NewManageSettings(settingsFile),
		Cleanup:      app.NewCleanupExpired(store, time.Now),
		ShowQR:       app.NewShowQR(clipboard.System{}, encoder.PNG{}, store, viewer.Window{}, time.Now),
		EditSettings: cli.EditSettingsInteractively,
		Version:      version,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}
	os.Exit(command.Run(os.Args[1:]))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "qr:", err)
	os.Exit(1)
}
