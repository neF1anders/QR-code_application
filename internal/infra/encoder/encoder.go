// Package encoder renders text as PNG QR codes using github.com/skip2/go-qrcode.
package encoder

import (
	"fmt"

	"github.com/skip2/go-qrcode"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

var recoveryLevels = map[domain.RecoveryLevel]qrcode.RecoveryLevel{
	domain.RecoveryLow:     qrcode.Low,
	domain.RecoveryMedium:  qrcode.Medium,
	domain.RecoveryHigh:    qrcode.High,
	domain.RecoveryHighest: qrcode.Highest,
}

// PNG encodes QR codes as square PNG images.
type PNG struct{}

// Encode returns a sizePx×sizePx PNG of text.
func (PNG) Encode(text string, sizePx int, recovery domain.RecoveryLevel) ([]byte, error) {
	level, ok := recoveryLevels[recovery]
	if !ok {
		return nil, fmt.Errorf("%w: %q", domain.ErrUnknownRecoveryLevel, recovery)
	}
	return qrcode.Encode(text, level, sizePx)
}
