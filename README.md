# qr

A tiny desktop CLI that shows text as a QR code in a window, so you can move a link from your PC to your phone without messengers or social networks: copy, type `qr`, scan.

The app is fully local: it never goes online. Only the install script downloads a release from GitHub, when you run it.

## Install / update

```sh
curl -fsSL https://raw.githubusercontent.com/neF1anders/qr/main/install.sh | sh
```

Running the same command again updates `qr` to the latest release. The binary goes to `~/.local/bin/qr` and needs no `sudo`; the script verifies the release checksum.

Options (environment variables): `QR_VERSION=v0.1.0` installs a specific release, `QR_INSTALL_DIR=/some/dir` changes the target directory.

Supported now: Linux x86_64 (tested on Ubuntu 24.04, GNOME on Wayland). Other systems will follow.

Uninstall: `rm ~/.local/bin/qr && rm -rf ~/.config/qr ~/.cache/qr`.

## Usage

```sh
qr                    # encode the clipboard text
qr https://example.com
qr some words here    # words are joined with spaces
qr -- --settings      # encode text that starts with "-"
qr --settings         # change settings with dropdown lists
qr --version
qr --help
```

The window closes with its close button or `Esc`.

## Settings

`qr --settings` saves to `~/.config/qr/config.json`:

```json
{
  "ttl_seconds": 0,
  "size_px": 384,
  "recovery_level": "medium"
}
```

| Setting | Meaning |
|---|---|
| `ttl_seconds` | How long the QR image stays on disk in `~/.cache/qr/`. `0` (default) deletes it as soon as the window closes. With `N > 0` the image is kept for N seconds; expired images are removed the next time `qr` runs. |
| `size_px` | Image size, 128–2048. |
| `recovery_level` | `low`, `medium`, `high` or `highest` error correction. Higher survives more damage, lower fits more text. |

A missing or broken config falls back to the defaults with a warning.

## Development

Requirements: Go (see `go.mod`) and the GUI headers used by [Fyne](https://fyne.io):

```sh
sudo apt-get install libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev
```

```sh
go run ./cmd/qr hello        # run
go test -race ./...          # tests
go vet ./...
gofmt -l .
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run
go run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean  # local release build
```

### Layout

Dependencies point inwards: `cli`/`infra` → `app` → `domain`.

```
cmd/qr/            composition root: wires adapters into use cases
internal/domain/   settings rules, QR image lifetime
internal/app/      use cases (show QR, clean up expired images, manage settings) and their ports
internal/infra/    clipboard, QR encoder, image storage, JSON config, Fyne window
internal/cli/      argument parsing and the interactive settings form
```

### Releasing

Push a tag `vX.Y.Z` to `main`. The `Release` workflow runs the tests, builds `qr_linux_amd64.tar.gz` with GoReleaser and publishes it with `checksums.txt` to GitHub Releases; `install.sh` picks it up from there.

## History

This started as a Go learning project: a global `Ctrl+F1` hotkey that served the QR code to a browser through `localhost`. It worked once per launch and kept running in the background, so it was rewritten as a CLI.
