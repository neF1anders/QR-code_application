package cli

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

var errFake = errors.New("fake failure")

type fakeSettings struct {
	current domain.Settings
	loadErr error
	saved   []domain.Settings
}

func (s *fakeSettings) Load() (domain.Settings, error) { return s.current, s.loadErr }

func (s *fakeSettings) Save(settings domain.Settings) error {
	s.saved = append(s.saved, settings)
	return nil
}

type fakeCleaner struct {
	ttls []time.Duration
	err  error
}

func (c *fakeCleaner) Execute(ttl time.Duration) error {
	c.ttls = append(c.ttls, ttl)
	return c.err
}

type fakeShower struct {
	fromClipboard int
	texts         []string
	err           error
}

func (s *fakeShower) FromClipboard(domain.Settings) error {
	s.fromClipboard++
	return s.err
}

func (s *fakeShower) FromText(text string, _ domain.Settings) error {
	s.texts = append(s.texts, text)
	return s.err
}

type fixture struct {
	cli      *CLI
	settings *fakeSettings
	cleaner  *fakeCleaner
	shower   *fakeShower
	edited   int
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
}

func newFixture() *fixture {
	f := &fixture{
		settings: &fakeSettings{current: domain.DefaultSettings()},
		cleaner:  &fakeCleaner{},
		shower:   &fakeShower{},
		stdout:   &bytes.Buffer{},
		stderr:   &bytes.Buffer{},
	}
	f.cli = &CLI{
		Settings: f.settings,
		Cleanup:  f.cleaner,
		ShowQR:   f.shower,
		EditSettings: func(current domain.Settings) (domain.Settings, error) {
			f.edited++
			current.TTL = time.Minute
			return current, nil
		},
		Version: "v1.2.3",
		Stdout:  f.stdout,
		Stderr:  f.stderr,
	}
	return f
}

func TestRunWithoutArgumentsEncodesClipboard(t *testing.T) {
	f := newFixture()

	if code := f.cli.Run(nil); code != exitOK {
		t.Fatalf("exit code = %d, stderr = %q", code, f.stderr)
	}
	if f.shower.fromClipboard != 1 || len(f.shower.texts) != 0 {
		t.Fatalf("clipboard calls = %d, text calls = %v", f.shower.fromClipboard, f.shower.texts)
	}
}

func TestRunWithArgumentsEncodesJoinedText(t *testing.T) {
	f := newFixture()

	if code := f.cli.Run([]string{"hello", "world"}); code != exitOK {
		t.Fatalf("exit code = %d", code)
	}
	if !slices.Equal(f.shower.texts, []string{"hello world"}) || f.shower.fromClipboard != 0 {
		t.Fatalf("texts = %v, clipboard calls = %d", f.shower.texts, f.shower.fromClipboard)
	}
}

func TestRunAfterDoubleDashEncodesFlagLikeText(t *testing.T) {
	f := newFixture()

	f.cli.Run([]string{"--", "--settings"})

	if !slices.Equal(f.shower.texts, []string{"--settings"}) || f.edited != 0 {
		t.Fatalf("texts = %v, edited = %d", f.shower.texts, f.edited)
	}
}

func TestRunCleansUpBeforeShowingWithConfiguredTTL(t *testing.T) {
	f := newFixture()
	f.settings.current.TTL = time.Hour

	f.cli.Run([]string{"x"})

	if !slices.Equal(f.cleaner.ttls, []time.Duration{time.Hour}) {
		t.Fatalf("cleanup ttls = %v", f.cleaner.ttls)
	}
}

func TestRunCleanupFailureIsOnlyAWarning(t *testing.T) {
	f := newFixture()
	f.cleaner.err = errFake

	if code := f.cli.Run([]string{"x"}); code != exitOK {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(f.stderr.String(), "warning") || len(f.shower.texts) != 1 {
		t.Fatalf("stderr = %q, texts = %v", f.stderr, f.shower.texts)
	}
}

func TestRunReportsShowFailure(t *testing.T) {
	f := newFixture()
	f.shower.err = errFake

	if code := f.cli.Run(nil); code != exitError {
		t.Fatalf("exit code = %d, want %d", code, exitError)
	}
	if !strings.Contains(f.stderr.String(), errFake.Error()) {
		t.Fatalf("stderr = %q", f.stderr)
	}
}

func TestRunSettingsSavesEditedSettings(t *testing.T) {
	f := newFixture()

	if code := f.cli.Run([]string{"--settings"}); code != exitOK {
		t.Fatalf("exit code = %d, stderr = %q", code, f.stderr)
	}
	if len(f.settings.saved) != 1 || f.settings.saved[0].TTL != time.Minute {
		t.Fatalf("saved = %+v", f.settings.saved)
	}
	if f.shower.fromClipboard != 0 {
		t.Fatal("--settings must not show a QR code")
	}
}

func TestRunSettingsCanceledDoesNotSave(t *testing.T) {
	f := newFixture()
	f.cli.EditSettings = func(domain.Settings) (domain.Settings, error) { return domain.Settings{}, ErrEditCanceled }

	if code := f.cli.Run([]string{"--settings"}); code != exitOK {
		t.Fatalf("exit code = %d", code)
	}
	if len(f.settings.saved) != 0 {
		t.Fatal("canceled edit must not be saved")
	}
}

func TestRunSettingsRejectsExtraArguments(t *testing.T) {
	f := newFixture()

	if code := f.cli.Run([]string{"--settings", "x"}); code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
}

func TestRunVersion(t *testing.T) {
	f := newFixture()

	if code := f.cli.Run([]string{"--version"}); code != exitOK || f.stdout.String() != "qr v1.2.3\n" {
		t.Fatalf("exit code = %d, stdout = %q", code, f.stdout)
	}
}

func TestRunHelpAndUnknownFlag(t *testing.T) {
	tests := []struct {
		args []string
		want int
	}{
		{[]string{"--help"}, exitOK},
		{[]string{"-h"}, exitOK},
		{[]string{"--bogus"}, exitUsage},
	}
	for _, tt := range tests {
		t.Run(tt.args[0], func(t *testing.T) {
			f := newFixture()
			if code := f.cli.Run(tt.args); code != tt.want {
				t.Fatalf("exit code = %d, want %d", code, tt.want)
			}
			if !strings.Contains(f.stderr.String(), "Usage:") {
				t.Fatalf("stderr = %q, want usage", f.stderr)
			}
		})
	}
}

func TestRunWarnsAndUsesDefaultsWhenSettingsUnreadable(t *testing.T) {
	f := newFixture()
	f.settings.loadErr = errFake

	if code := f.cli.Run([]string{"x"}); code != exitOK {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(f.stderr.String(), "warning") {
		t.Fatalf("stderr = %q", f.stderr)
	}
}
