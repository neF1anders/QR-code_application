// Package viewer shows QR images in a native Fyne window.
package viewer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const appID = "io.github.nef1anders.qr"

// Window opens one window per Show call.
type Window struct{}

// Show displays png with the encoded text as a caption and returns once the
// window is closed (close button or Esc).
func (Window) Show(png []byte, caption string) error {
	application := app.NewWithID(appID)
	window := application.NewWindow("qr")

	image := canvas.NewImageFromResource(fyne.NewStaticResource("qr.png", png))
	image.FillMode = canvas.ImageFillContain
	// Nearest-neighbour keeps module edges sharp, which scanners need.
	image.ScaleMode = canvas.ImageScalePixels
	image.SetMinSize(fyne.NewSize(320, 320))

	label := widget.NewLabel(caption)
	label.Alignment = fyne.TextAlignCenter
	label.Truncation = fyne.TextTruncateEllipsis

	window.SetContent(container.NewBorder(nil, label, nil, nil, image))
	window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEscape {
			window.Close()
		}
	})
	window.CenterOnScreen()
	window.ShowAndRun()
	return nil
}
