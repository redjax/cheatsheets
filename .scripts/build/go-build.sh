#!/usr/bin/env bash
set -euo pipefail

## Get to project root (script is in .scripts/build/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Building chtsht with Go"

## Clean up dist directory
if [[ -d "$PROJECT_ROOT/dist" ]]; then
    echo "Removing existing dist/ directory"
    rm -rf "$PROJECT_ROOT/dist"
fi

## Create dist directory
mkdir -p "$PROJECT_ROOT/dist"

## Build for current platform
cd "$PROJECT_ROOT/app"
go build -o "$PROJECT_ROOT/dist/chtsht" -trimpath -ldflags="-s -w" ./cmd/chtsht/main.go

echo "Built binary to dist/chtsht"
