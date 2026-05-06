#!/usr/bin/env bash

############################################################################
# When cloning this repo, it defaults to the https URL; often times I      #
# already have SSH keys added to the remote, and so I want to use SSH for  #
# git operations.                                                          #
#                                                                          #
# This script is a helper to switch the origin between https:// to ssh://. #
############################################################################

set -uo pipefail

if ! command -v git &>/dev/null; then
  echo "[ERROR] git is not installed."
  exit 1
fi

##
# Default values
##

## SSH url
REMOTE_SSH_URL="git@github.com:redjax/cheatsheets.git"
## HTTPS url
REMOTE_HTTPS_URL="https://github.com/redjax/cheatsheets.git"

## Which URL to use (https/ssh)
REMOTE_TYPE="ssh"
## Empty URL value to populate later
REMOTE_URL=""

function usage() {
  echo ""
  echo "Usage: ${0} [OPTIONS]"
  echo ""
  echo "Options:"
  echo "  -h|--help          Print this help menu."
  echo "  --ssh-url          Remote repository's SSH URL, i.e. git@github.com/redjax/cheatsheets.git"
  echo "  --https-url        Remote repository's HTTPS URL, i.e. https://github.com/redjax/cheatsheets"
  echo "  -t, --remote-type  Either 'https' or 'ssh'"
  echo ""
}

## Parse args
while [[ $# -gt 0 ]]; do
  case $1 in
  --ssh-url)
    if [[ -z "${2:-}" ]]; then
      echo "[ERROR] --ssh-url provided, but no URL given."

      usage
      exit 1
    fi

    REMOTE_SSH_URL="$2"
    shift 2
    ;;
  --https-url)
    if [[ -z "${2:-}" ]]; then
      echo "[ERROR] --https-url provided, but no URL given."

      usage
      exit 1
    fi

    REMOTE_HTTPS_URL="$2"
    shift 2
    ;;
  -t | --remote-type)
    if [[ -z "${2:-}" ]]; then
      echo "[ERROR] --remote-type provided, but no 'https'/'ssh' value given."

      usage
      exit 1
    fi

    ## Ensure value is https or ssh
    value="$(echo "$2" | tr -d '[:space:]' | tr '[:upper:]' '[:lower:]')"

    case "$value" in
    https)
      REMOTE_TYPE="https"
      ;;
    ssh)
      REMOTE_TYPE="ssh"
      ;;
    *)
      echo "[ERROR] Invalid --remote-type: $2. Valid values are: https, ssh"
      usage
      exit 1
      ;;
    esac

    shift 2
    ;;
  -h | -help)
    usage
    exit 0
    ;;
  *)
    echo "[ERROR] Invalid argument: $1"

    usage
    exit 1
    ;;
  esac
done

if [[ "${REMOTE_TYPE}" == "https" ]]; then
  REMOTE_URL="${REMOTE_HTTPS_URL}"
else
  REMOTE_URL="${REMOTE_SSH_URL}"
fi

## Remove existing origin
echo "Removing current origin if it exists"
if ! git remote remove origin; then
  echo "[WARNING] Failed removing existing remote."
  echo "  Continuing anyway, in case there was no origin to begin with."
fi

## Add new origin
echo "Setting dotfiles remote repository URL to: (${REMOTE_TYPE}) ${REMOTE_URL}"
if ! git remote add origin "$REMOTE_URL"; then
  echo "[ERROR] Failed to add new origin."
  exit 1
fi

## Fetch remote branches so tracking info can be updated
echo "Fetching remote branches to update local upstream URLs"
git fetch origin --prune

## Try to restore upstream tracking for each local branch
echo "Updating local branches to track remote branches (if available)"
for branch in $(git for-each-ref --format='%(refname:short)' refs/heads/); do
  ## Check if branch exists on remote
  if git show-ref --verify --quiet "refs/remotes/origin/$branch"; then
    ## Update branch
    git branch --set-upstream-to="origin/$branch" "$branch" 2>/dev/null &&
      echo "[OK] Linked local branch '$branch' -> origin/$branch" ||
      echo "[WARN] Failed to set upstream for branch '$branch'"
  else
    echo "[INFO] Skipping '$branch' (not on remote)"
  fi
done

echo ""
echo "[SUCCESS] Set remote origin URL to: (${REMOTE_TYPE}) ${REMOTE_URL}"
exit 0

