// Package storage keeps rendered QR images as PNG files in one directory.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/neF1anders/qr/internal/domain"
)

// The creation time is encoded in the file name so lifetimes do not depend on file mtimes.
const fileExtension = ".png"

// Directory stores images in a single directory, created on first save.
type Directory struct {
	path string
}

func NewDirectory(path string) *Directory {
	return &Directory{path: path}
}

func (d *Directory) Save(png []byte, createdAt time.Time) (domain.StoredQR, error) {
	if err := os.MkdirAll(d.path, 0o700); err != nil {
		return domain.StoredQR{}, err
	}
	name := strconv.FormatInt(createdAt.UnixNano(), 10) + fileExtension
	image := domain.StoredQR{Path: filepath.Join(d.path, name), CreatedAt: createdAt}
	if err := os.WriteFile(image.Path, png, 0o600); err != nil {
		return domain.StoredQR{}, err
	}
	return image, nil
}

// List returns the stored images; files not created by Save are ignored.
func (d *Directory) List() ([]domain.StoredQR, error) {
	entries, err := os.ReadDir(d.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var images []domain.StoredQR
	for _, entry := range entries {
		createdAt, ok := parseCreatedAt(entry.Name())
		if !ok || entry.IsDir() {
			continue
		}
		images = append(images, domain.StoredQR{Path: filepath.Join(d.path, entry.Name()), CreatedAt: createdAt})
	}
	return images, nil
}

// Delete removes the image; an already missing file is not an error.
func (d *Directory) Delete(image domain.StoredQR) error {
	if filepath.Dir(image.Path) != filepath.Clean(d.path) {
		return fmt.Errorf("refusing to delete %s outside %s", image.Path, d.path)
	}
	if err := os.Remove(image.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func parseCreatedAt(fileName string) (time.Time, bool) {
	stem, found := strings.CutSuffix(fileName, fileExtension)
	if !found {
		return time.Time{}, false
	}
	nanos, err := strconv.ParseInt(stem, 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(0, nanos), true
}
