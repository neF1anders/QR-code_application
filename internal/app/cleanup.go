package app

import (
	"errors"
	"fmt"
	"time"
)

// CleanupExpired deletes stored QR images that have outlived the configured TTL.
type CleanupExpired struct {
	store ImageStore
	now   Clock
}

func NewCleanupExpired(store ImageStore, now Clock) *CleanupExpired {
	return &CleanupExpired{store: store, now: now}
}

// Execute removes every expired image and reports all deletions that failed.
func (u *CleanupExpired) Execute(ttl time.Duration) error {
	images, err := u.store.List()
	if err != nil {
		return fmt.Errorf("list qr images: %w", err)
	}
	now := u.now()
	var errs []error
	for _, image := range images {
		if !image.IsExpired(now, ttl) {
			continue
		}
		if err := u.store.Delete(image); err != nil {
			errs = append(errs, fmt.Errorf("delete %s: %w", image.Path, err))
		}
	}
	return errors.Join(errs...)
}
