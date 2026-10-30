package domain

import "time"

// StoredQR is a rendered QR image kept on disk.
type StoredQR struct {
	Path      string
	CreatedAt time.Time
}

// IsExpired reports whether the image has outlived ttl at the moment now.
// With a zero ttl every stored image is expired: it should not survive its window.
func (q StoredQR) IsExpired(now time.Time, ttl time.Duration) bool {
	return !now.Before(q.CreatedAt.Add(ttl))
}
