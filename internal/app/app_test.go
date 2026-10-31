package app

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/neF1anders/qr/internal/domain"
)

type showQRFixture struct {
	encoder *fakeEncoder
	store   *fakeStore
	viewer  *fakeViewer
}

func newShowQR(clipboard fakeClipboard) (*ShowQR, showQRFixture) {
	f := showQRFixture{encoder: &fakeEncoder{}, store: &fakeStore{}, viewer: &fakeViewer{}}
	return NewShowQR(clipboard, f.encoder, f.store, f.viewer, fixedClock), f
}

func TestShowQRFromTextWithZeroTTLDeletesAfterWindowCloses(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{})

	if err := useCase.FromText("hello", domain.DefaultSettings()); err != nil {
		t.Fatalf("FromText() error = %v", err)
	}
	if !slices.Equal(f.viewer.shownCaptions, []string{"hello"}) {
		t.Fatalf("shown = %v, want [hello]", f.viewer.shownCaptions)
	}
	if len(f.store.deleted) != 1 {
		t.Fatalf("deleted %d images, want 1", len(f.store.deleted))
	}
}

func TestShowQRFromTextWithTTLKeepsImage(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{})
	settings := domain.DefaultSettings()
	settings.TTL = time.Minute

	if err := useCase.FromText("hello", settings); err != nil {
		t.Fatalf("FromText() error = %v", err)
	}
	if len(f.store.deleted) != 0 {
		t.Fatalf("deleted %d images, want 0", len(f.store.deleted))
	}
}

func TestShowQRFromTextRejectsBlankText(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{})

	err := useCase.FromText(" \n\t", domain.DefaultSettings())
	if !errors.Is(err, ErrNothingToEncode) {
		t.Fatalf("FromText() error = %v, want ErrNothingToEncode", err)
	}
	if len(f.encoder.encoded) != 0 {
		t.Fatal("blank text must not be encoded")
	}
}

func TestShowQRFromClipboardEncodesClipboardText(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{text: "https://example.com"})

	if err := useCase.FromClipboard(domain.DefaultSettings()); err != nil {
		t.Fatalf("FromClipboard() error = %v", err)
	}
	if !slices.Equal(f.encoder.encoded, []string{"https://example.com"}) {
		t.Fatalf("encoded = %v", f.encoder.encoded)
	}
}

func TestShowQRFromClipboardPropagatesReadError(t *testing.T) {
	useCase, _ := newShowQR(fakeClipboard{err: errFake})

	if err := useCase.FromClipboard(domain.DefaultSettings()); !errors.Is(err, errFake) {
		t.Fatalf("FromClipboard() error = %v, want errFake", err)
	}
}

func TestShowQRDeletesImageEvenWhenViewerFails(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{})
	f.viewer.err = errFake

	if err := useCase.FromText("hello", domain.DefaultSettings()); !errors.Is(err, errFake) {
		t.Fatalf("FromText() error = %v, want errFake", err)
	}
	if len(f.store.deleted) != 1 {
		t.Fatalf("deleted %d images, want 1", len(f.store.deleted))
	}
}

func TestShowQRPropagatesEncoderError(t *testing.T) {
	useCase, f := newShowQR(fakeClipboard{})
	f.encoder.err = errFake

	if err := useCase.FromText("hello", domain.DefaultSettings()); !errors.Is(err, errFake) {
		t.Fatalf("FromText() error = %v, want errFake", err)
	}
	if len(f.viewer.shownCaptions) != 0 {
		t.Fatal("viewer must not open when encoding fails")
	}
}

func TestCleanupExpiredDeletesOnlyExpiredImages(t *testing.T) {
	fresh := domain.StoredQR{Path: "fresh", CreatedAt: fixedNow.Add(-30 * time.Second)}
	old := domain.StoredQR{Path: "old", CreatedAt: fixedNow.Add(-2 * time.Minute)}
	store := &fakeStore{images: []domain.StoredQR{fresh, old}}

	if err := NewCleanupExpired(store, fixedClock).Execute(time.Minute); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !slices.Equal(store.deleted, []domain.StoredQR{old}) {
		t.Fatalf("deleted = %v, want [old]", store.deleted)
	}
}

func TestCleanupExpiredReportsDeleteFailures(t *testing.T) {
	store := &fakeStore{
		images:    []domain.StoredQR{{Path: "old", CreatedAt: fixedNow.Add(-time.Hour)}},
		deleteErr: errFake,
	}

	if err := NewCleanupExpired(store, fixedClock).Execute(time.Minute); !errors.Is(err, errFake) {
		t.Fatalf("Execute() error = %v, want errFake", err)
	}
}

func TestManageSettingsLoadFallsBackToDefaults(t *testing.T) {
	invalid := domain.DefaultSettings()
	invalid.SizePx = 1
	tests := []struct {
		name       string
		repository *fakeSettingsRepository
	}{
		{"unreadable", &fakeSettingsRepository{loadErr: errFake}},
		{"invalid", &fakeSettingsRepository{settings: invalid}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings, err := NewManageSettings(tt.repository).Load()
			if err == nil {
				t.Fatal("Load() must report why defaults are used")
			}
			if settings != domain.DefaultSettings() {
				t.Fatalf("Load() = %+v, want defaults", settings)
			}
		})
	}
}

func TestManageSettingsSaveRejectsInvalidSettings(t *testing.T) {
	repository := &fakeSettingsRepository{}
	invalid := domain.DefaultSettings()
	invalid.TTL = -time.Second

	if err := NewManageSettings(repository).Save(invalid); !errors.Is(err, domain.ErrNegativeTTL) {
		t.Fatalf("Save() error = %v, want ErrNegativeTTL", err)
	}
	if len(repository.saved) != 0 {
		t.Fatal("invalid settings must not be persisted")
	}
}
