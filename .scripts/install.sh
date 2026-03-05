#!/usr/bin/env bash
#
# Install script for chtsht (cheatsheets CLI)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/redjax/cheatsheets/main/.scripts/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/redjax/cheatsheets/main/.scripts/install.sh | bash -s -- --auto
#

set -e

REPO="redjax/cheatsheets"
BIN_NAME="chtsht"
INSTALL_DIR="$HOME/.local/bin"
AUTO_MODE=false

## Parse arguments
for arg in "$@"; do
  case "$arg" in
    --auto) AUTO_MODE=true ;;
    --help|-h)
      echo "Usage: install.sh [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --auto    Install without prompts"
      echo "  --help    Show this help message"
      exit 0
      ;;
  esac
done

## Check dependencies
for cmd in curl unzip; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "[ERROR] '$cmd' is required but not installed."
    exit 1
  fi
done

## Check if already installed
if command -v "$BIN_NAME" &>/dev/null; then
  EXISTING_PATH=$(command -v "$BIN_NAME")
  if [[ "$AUTO_MODE" == true ]]; then
    echo "$BIN_NAME is already installed at $EXISTING_PATH. Reinstalling"
  else
    read -p "$BIN_NAME is already installed at $EXISTING_PATH. Download and install anyway? (y/N) " CONFIRM
    if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
      echo "Aborting."
      exit 0
    fi
  fi
fi

## Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux)  OS_ASSET="linux" ;;
  Darwin) OS_ASSET="macOS" ;;
  *)
    echo "[ERROR] Unsupported OS: $OS"
    exit 1
    ;;
esac

## Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)         ARCH_ASSET="amd64" ;;
  aarch64|arm64)  ARCH_ASSET="arm64" ;;
  *)
    echo "[ERROR] Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

## Get latest release version from GitHub API
echo "Fetching latest release info"
RELEASE_TAG=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep -Po '"tag_name": "\K.*?(?=")')

if [[ -z "$RELEASE_TAG" ]]; then
  echo "[ERROR] Could not determine latest release version."
  exit 1
fi

VERSION="${RELEASE_TAG#v}"
echo "Latest version: $RELEASE_TAG"

## Build asset name: chtsht-linux-amd64-0.1.0.zip
FILE="${BIN_NAME}-${OS_ASSET}-${ARCH_ASSET}-${VERSION}.zip"
URL="https://github.com/$REPO/releases/download/${RELEASE_TAG}/${FILE}"

## Create temp directory
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

ARCHIVE="$TMPDIR/$FILE"

## Download
echo "Downloading $FILE"
if ! curl -fSL -o "$ARCHIVE" "$URL"; then
  echo "[ERROR] Download failed. Check that a release exists for $OS_ASSET/$ARCH_ASSET."
  echo "  URL: $URL"
  exit 1
fi

## Extract
echo "Extracting"
unzip -q "$ARCHIVE" -d "$TMPDIR"

## Verify the binary exists in the extracted files
BINARY_PATH="$TMPDIR/$BIN_NAME"
if [[ ! -f "$BINARY_PATH" ]]; then
  echo "[ERROR] Binary '$BIN_NAME' not found in archive."
  echo "Archive contents:"
  ls -la "$TMPDIR"
  exit 1
fi

## Ensure install directory exists
mkdir -p "$INSTALL_DIR"

## Install
chmod 755 "$BINARY_PATH"
if [[ -w "$INSTALL_DIR" ]]; then
  mv "$BINARY_PATH" "$INSTALL_DIR/$BIN_NAME"
else
  echo "Install directory $INSTALL_DIR is not writable, using sudo"
  sudo mv "$BINARY_PATH" "$INSTALL_DIR/$BIN_NAME"
fi

echo ""
echo "$BIN_NAME $RELEASE_TAG installed to $INSTALL_DIR/$BIN_NAME"

## Check if install dir is in PATH
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo ""
  echo "WARNING:  $INSTALL_DIR is not in your PATH."
  echo ""
  echo "Add it by appending this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
  echo ""
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  echo ""
  echo "Then restart your shell or run: source ~/.bashrc"
fi

## Verify installation
if command -v "$BIN_NAME" &>/dev/null; then
  echo ""
  "$BIN_NAME" self version
fi
