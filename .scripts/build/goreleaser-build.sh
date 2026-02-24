#!/usr/bin/env bash
set -euo pipefail

if ! command -v goreleaser &> /dev/null; then
  echo "[ERROR] goreleaser is not installed. Please install it first." >&2
  exit 1
fi

LOCAL="false"

function usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help  Show this help message and exit"
    echo "  --local     Build locally without creating a Git tag or pushing to remote"
    echo ""
}

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --local)
        LOCAL="true"
        shift
        ;;
    -h | --help)
        usage
        exit 0
        ;;
    *)
        echo "Invalid arg: $1" >&2
        usage
        exit 1
        ;;
  esac
done

echo "Building chtsht with goreleaser"

if [[ "$LOCAL" == "true" ]]; then
    echo "Doing local goreleaser build (no Git tag, no push)"
    if ! goreleaser build --clean --snapshot; then
        echo "[ERROR] goreleaser local build failed." >&2
        exit 1
    fi
    
    ## Detect current platform and copy binary to root
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    ## Normalize OS name
    if [[ "$OS" == "darwin" ]]; then
        OS_NAME="darwin"
    else
        OS_NAME="linux"
    fi
    
    ## Normalize architecture
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH_NAME="amd64"
    elif [[ "$ARCH" == "aarch64" ]] || [[ "$ARCH" == "arm64" ]]; then
        ARCH_NAME="arm64"
    else
        echo "[ERROR] Unsupported architecture: $ARCH" >&2
        exit 1
    fi
    
    ## Find the binary in dist/
    BINARY_DIR="dist/chtsht_${OS_NAME}_${ARCH_NAME}"*
    BINARY_PATH=$(find $BINARY_DIR -name "chtsht" -type f 2>/dev/null | head -n 1)
    
    if [[ -z "$BINARY_PATH" ]]; then
        echo "[ERROR] Could not find binary for ${OS_NAME}_${ARCH_NAME} in dist/" >&2
        exit 1
    fi
    
    ## Copy binary to root
    cp "$BINARY_PATH" ./chtsht
    chmod +x ./chtsht
    echo "Copied binary to ./chtsht (${OS_NAME}_${ARCH_NAME})"
else
    if ! goreleaser release --clean; then
        echo "[ERROR] goreleaser build failed." >&2
        exit 1
    fi
fi
