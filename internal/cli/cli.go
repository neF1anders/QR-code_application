// Package cli is the only entry point of the qr utility: it parses arguments and calls use cases.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/neF1anders/qr/internal/domain"
)

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

const usage = `qr - show text as a QR code in a window

Usage:
  qr                 encode the clipboard text
  qr <text...>       encode the given text (words are joined with spaces)
  qr -- <text...>    encode text that starts with "-"
  qr --settings      change settings interactively
  qr --version       print the version
  qr --help          show this help
`

type settingsService interface {
	Load() (domain.Settings, error)
	Save(settings domain.Settings) error
}

type imageCleaner interface {
	Execute(ttl time.Duration) error
}

type qrShower interface {
	FromClipboard(settings domain.Settings) error
	FromText(text string, settings domain.Settings) error
}

// SettingsEditor asks the user for new settings, starting from current.
type SettingsEditor func(current domain.Settings) (domain.Settings, error)

// ErrEditCanceled is returned by a SettingsEditor when the user leaves without saving.
var ErrEditCanceled = errors.New("settings edit canceled")

// CLI wires command-line arguments to the use cases.
type CLI struct {
	Settings     settingsService
	Cleanup      imageCleaner
	ShowQR       qrShower
	EditSettings SettingsEditor
	Version      string
	Stdout       io.Writer
	Stderr       io.Writer
}

// Run executes the command described by args (without the program name) and returns the exit code.
func (c *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("qr", flag.ContinueOnError)
	flags.SetOutput(c.Stderr)
	flags.Usage = func() { fmt.Fprint(c.Stderr, usage) }
	editSettings := flags.Bool("settings", false, "change settings interactively")
	printVersion := flags.Bool("version", false, "print the version")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}
		return exitUsage
	}

	switch {
	case *printVersion:
		fmt.Fprintln(c.Stdout, "qr", c.Version)
		return exitOK
	case *editSettings && flags.NArg() > 0:
		fmt.Fprintln(c.Stderr, "qr: --settings takes no arguments")
		return exitUsage
	case *editSettings:
		return c.runSettings()
	default:
		return c.runShow(flags.Args())
	}
}

func (c *CLI) runSettings() int {
	current := c.loadSettings()
	updated, err := c.EditSettings(current)
	if errors.Is(err, ErrEditCanceled) {
		fmt.Fprintln(c.Stdout, "Settings unchanged.")
		return exitOK
	}
	if err != nil {
		return c.fail(err)
	}
	if err := c.Settings.Save(updated); err != nil {
		return c.fail(err)
	}
	fmt.Fprintln(c.Stdout, "Settings saved.")
	return exitOK
}

func (c *CLI) runShow(words []string) int {
	settings := c.loadSettings()
	if err := c.Cleanup.Execute(settings.TTL); err != nil {
		c.warn(err)
	}

	var err error
	if len(words) == 0 {
		err = c.ShowQR.FromClipboard(settings)
	} else {
		err = c.ShowQR.FromText(strings.Join(words, " "), settings)
	}
	if err != nil {
		return c.fail(err)
	}
	return exitOK
}

// loadSettings never fails: unreadable settings are reported and replaced by defaults.
func (c *CLI) loadSettings() domain.Settings {
	settings, err := c.Settings.Load()
	if err != nil {
		c.warn(err)
	}
	return settings
}

func (c *CLI) warn(err error) {
	fmt.Fprintln(c.Stderr, "qr: warning:", err)
}

func (c *CLI) fail(err error) int {
	fmt.Fprintln(c.Stderr, "qr:", err)
	return exitError
}
