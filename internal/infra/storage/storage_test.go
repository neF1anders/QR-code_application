package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

func TestSaveListDeleteRoundTrip(t *testing.T) {
	dir := NewDirectory(filepath.Join(t.TempDir(), "qr"))
	createdAt := time.Unix(1_700_000_000, 42)

	saved, err := dir.Save([]byte("png"), createdAt)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	images, err := dir.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(images) != 1 || images[0].Path != saved.Path || !images[0].CreatedAt.Equal(createdAt) {
		t.Fatalf("List() = %+v, want [%+v]", images, saved)
	}

	if err := dir.Delete(saved); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := os.Stat(saved.Path); !os.IsNotExist(err) {
		t.Fatalf("file still exists: %v", err)
	}
	if err := dir.Delete(saved); err != nil {
		t.Fatalf("second Delete() error = %v, want nil", err)
	}
}

func TestListOnMissingDirectoryIsEmpty(t *testing.T) {
	images, err := NewDirectory(filepath.Join(t.TempDir(), "absent")).List()
	if err != nil || len(images) != 0 {
		t.Fatalf("List() = %v, %v; want empty, nil", images, err)
	}
}

func TestListIgnoresForeignFiles(t *testing.T) {
	path := t.TempDir()
	for _, name := range []string{"notes.txt", "photo.png", "123.jpg"} {
		if err := os.WriteFile(filepath.Join(path, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	images, err := NewDirectory(path).List()
	if err != nil || len(images) != 0 {
		t.Fatalf("List() = %v, %v; want empty, nil", images, err)
	}
}

func TestDeleteRefusesPathsOutsideDirectory(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "keep.png")
	if err := os.WriteFile(outside, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := NewDirectory(t.TempDir()).Delete(domain.StoredQR{Path: outside}); err == nil {
		t.Fatal("Delete() must refuse a file outside its directory")
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatalf("outside file was removed: %v", err)
	}
}
