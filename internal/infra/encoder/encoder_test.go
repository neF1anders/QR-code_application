package encoder

import (
	"bytes"
	"errors"
	"image/png"
	"strings"
	"testing"

	"github.com/neF1anders/qr/internal/domain"
)

func TestEncodeProducesSquarePNGOfRequestedSize(t *testing.T) {
	for _, level := range domain.RecoveryLevels {
		t.Run(string(level), func(t *testing.T) {
			data, err := PNG{}.Encode("https://example.com", 256, level)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("output is not a PNG: %v", err)
			}
			if b := img.Bounds(); b.Dx() != 256 || b.Dy() != 256 {
				t.Fatalf("size = %dx%d, want 256x256", b.Dx(), b.Dy())
			}
		})
	}
}

func TestEncodeRejectsUnknownRecoveryLevel(t *testing.T) {
	_, err := PNG{}.Encode("x", 256, "ultra")
	if !errors.Is(err, domain.ErrUnknownRecoveryLevel) {
		t.Fatalf("Encode() error = %v, want ErrUnknownRecoveryLevel", err)
	}
}

func TestEncodeRejectsTextTooLongForQR(t *testing.T) {
	if _, err := (PNG{}).Encode(strings.Repeat("a", 5000), 256, domain.RecoveryHighest); err == nil {
		t.Fatal("Encode() must fail for text that does not fit a QR code")
	}
}
