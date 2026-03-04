#!/usr/bin/env bash
set -euo pipefail

if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go to run the tests."
    exit 1
fi

THIS_DIR=$(dirname "$(realpath "$0")")
PROJECT_ROOT=$(realpath -m "$THIS_DIR/..")
APP_DIR="${PROJECT_ROOT}/app"
ORIGINAL_DIR=$(pwd)
trap "cd \"$ORIGINAL_DIR\"" EXIT

echo "Running Go tests"
cd "$APP_DIR"

if ! go test -v ./...; then
    echo "Go tests failed"
    exit 1
fi
