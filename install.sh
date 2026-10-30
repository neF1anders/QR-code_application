#!/bin/sh
# Installs or updates the qr utility from GitHub Releases.
#
#   curl -fsSL https://raw.githubusercontent.com/neF1anders/QR-code_application/main/install.sh | sh
#
# Environment:
#   QR_VERSION      release tag to install, e.g. v0.1.0 (default: latest)
#   QR_INSTALL_DIR  target directory (default: ~/.local/bin)
#   QR_RELEASES_URL base URL of the releases (default: the GitHub repository)
set -eu

repo="neF1anders/QR-code_application"
releases_url="${QR_RELEASES_URL:-https://github.com/$repo/releases}"
version="${QR_VERSION:-latest}"
install_dir="${QR_INSTALL_DIR:-$HOME/.local/bin}"
archive="qr_linux_amd64.tar.gz"

fail() {
	echo "qr install: $*" >&2
	exit 1
}

[ "$(uname -s)" = "Linux" ] || fail "only Linux is supported for now"
case "$(uname -m)" in
x86_64 | amd64) ;;
*) fail "only x86_64 is supported for now, got $(uname -m)" ;;
esac
command -v curl >/dev/null || fail "curl is required"
command -v sha256sum >/dev/null || fail "sha256sum is required"

if [ "$version" = "latest" ]; then
	download_url="$releases_url/latest/download"
else
	download_url="$releases_url/download/$version"
fi

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT INT TERM

echo "Downloading qr ($version)..."
curl -fsSL -o "$workdir/$archive" "$download_url/$archive" || fail "download failed: $download_url/$archive"
curl -fsSL -o "$workdir/checksums.txt" "$download_url/checksums.txt" || fail "download failed: $download_url/checksums.txt"

(cd "$workdir" && grep " $archive\$" checksums.txt | sha256sum -c --status) || fail "checksum mismatch, aborting"

tar -xzf "$workdir/$archive" -C "$workdir" qr

old_version=""
if [ -x "$install_dir/qr" ]; then
	old_version="$("$install_dir/qr" --version 2>/dev/null || true)"
fi

mkdir -p "$install_dir"
# Install via rename so a running qr keeps its old binary intact.
cp "$workdir/qr" "$install_dir/.qr.new"
chmod 755 "$install_dir/.qr.new"
mv -f "$install_dir/.qr.new" "$install_dir/qr"

new_version="$("$install_dir/qr" --version)"
if [ "$old_version" = "$new_version" ]; then
	echo "Already up to date: $new_version"
elif [ -n "$old_version" ]; then
	echo "Updated: $old_version -> $new_version"
else
	echo "Installed: $new_version at $install_dir/qr"
fi

case ":$PATH:" in
*":$install_dir:"*) ;;
*) echo "Note: $install_dir is not in PATH. Add it, e.g.: echo 'export PATH=\"\$PATH:$install_dir\"' >> ~/.bashrc" ;;
esac
