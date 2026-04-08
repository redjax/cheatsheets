#!/usr/bin/bash
set -euo pipefail

THIS_DIR="$(dirname "${0}")"
REPO_ROOT=$(realpath -m $"{THIS_DIR}/..")
APP_ROOT="${REPO_ROOT}/app"
CWD=$(pwd)

if ! command -v go &>/dev/null; then
  echo "[ERROR] go is not installed"
  exit 1
fi

function cleanup() {
  cd "${CWD}"
}
trap cleanup EXIT

cd "${APP_ROOT}"
echo "Bumping packages in go.mod/go.sum"

if ! go get -u ./...; then
  echo "[ERROR] Failed bumping packages"
  exit 1
fi

