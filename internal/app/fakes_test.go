package app

import (
	"errors"
	"time"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

var errFake = errors.New("fake failure")

var fixedNow = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedNow }

type fakeClipboard struct {
	text string
	err  error
}

func (c fakeClipboard) ReadText() (string, error) { return c.text, c.err }

type fakeEncoder struct {
	encoded []string
	err     error
}

func (e *fakeEncoder) Encode(text string, _ int, _ domain.RecoveryLevel) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	e.encoded = append(e.encoded, text)
	return []byte("png:" + text), nil
}

type fakeStore struct {
	images    []domain.StoredQR
	deleted   []domain.StoredQR
	deleteErr error
}

func (s *fakeStore) Save(_ []byte, createdAt time.Time) (domain.StoredQR, error) {
	image := domain.StoredQR{Path: createdAt.String(), CreatedAt: createdAt}
	s.images = append(s.images, image)
	return image, nil
}

func (s *fakeStore) List() ([]domain.StoredQR, error) { return s.images, nil }

func (s *fakeStore) Delete(qr domain.StoredQR) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	s.deleted = append(s.deleted, qr)
	return nil
}

type fakeViewer struct {
	shownCaptions []string
	err           error
}

func (v *fakeViewer) Show(_ []byte, caption string) error {
	v.shownCaptions = append(v.shownCaptions, caption)
	return v.err
}

type fakeSettingsRepository struct {
	settings domain.Settings
	saved    []domain.Settings
	loadErr  error
}

func (r *fakeSettingsRepository) Load() (domain.Settings, error) { return r.settings, r.loadErr }

func (r *fakeSettingsRepository) Save(settings domain.Settings) error {
	r.saved = append(r.saved, settings)
	return nil
}
