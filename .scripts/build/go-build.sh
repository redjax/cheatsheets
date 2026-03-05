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

## Collect build metadata
GIT_VERSION=$(git -C "$PROJECT_ROOT" describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(git -C "$PROJECT_ROOT" rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

echo "Version: $GIT_VERSION  Commit: $GIT_COMMIT  Date: $BUILD_DATE"

## Build for current platform with version info injected
cd "$PROJECT_ROOT/app"
go build -o "$PROJECT_ROOT/dist/chtsht" -trimpath \
  -ldflags="-s -w \
    -X 'github.com/redjax/cheatsheets/internal/version.Version=${GIT_VERSION}' \
    -X 'github.com/redjax/cheatsheets/internal/version.Commit=${GIT_COMMIT}' \
    -X 'github.com/redjax/cheatsheets/internal/version.Date=${BUILD_DATE}'" \
  ./cmd/chtsht/main.go

echo "Built binary to dist/chtsht"
