---
description: "Linux OS."
last_updated: "2026-05-01"
tags: ["linux", "system", "os"]
---

# Linux <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [Paths](#paths)
  - [Variables](#variables)
- [Args](#args)
  - [Example: Parse CLI args](#example-parse-cli-args)
- [Bash Tips & Tricks](#bash-tips-tricks)
  - [Get command that launched a PID](#get-command-that-launched-a-pid)
  - [Only run function if script is called directly](#only-run-function-if-script-is-called-directly)
  - [Uppercase a variable's text value](#uppercase-a-variables-text-value)
- [Traps](#traps)
  - [Return to PWD on script exit](#return-to-pwd-on-script-exit)

## Paths

### Variables

| Variable                                                   | Description                             |
| ---------------------------------------------------------- | --------------------------------------- |
| `THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"` | Get parent directory of current script. |
| `REL_PATH=$(realpath -m "${THIS_DIR}/../..")`              | Get absolute path from a relative path. |

## Args

### Example: Parse CLI args

```shell
#!/usr/bin/env bash
set -uo pipefail

## Arg default values
DEBUG=0
NAME="world"

## Parse arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --debug|-D)
            DEBUG=1
            ;;
        -n|--name)
          if [[ -z "$2" ]]; then
            echo "[ERROR] --name provided, but no name given."
            exit 1
          fi

          NAME="$2"
          shift 2
          ;;
        --foo)
            FOO=1
            ;;
        --bar)
            BAR=1
            ;;
        ## catch-all for unknown options
        --*)
            echo "Warning: Unknown option $1"
            ;;
        ## Unmatched/unhandled
        *)
            ## Handle positional args
            #  You have to implement this based on the needs of the current script
            ;;
    esac
    shift
done

```

## Bash Tips & Tricks

### Get command that launched a PID

To get the command that launched a process, find its PID using `top`/`htop`/`btop` and run the following command, replacing `$PID` with your process ID:

```shell
tf '\0' ' ' < /proc/$PID/cmdline
```

For example, to find the command that started process `4645`:

```shell
tf '\0' ' ' < /proc/4645/cmdline
```

### Only run function if script is called directly

If you want a Bash script to only run somethin when it's called directly, i.e. `./path/to/script.sh --arg example`, you can add the following to the bottom of the script; the `main` function will only run when the script is caled directy, meaning you can source other scripts to import their functions/variables (`i.e. . /path/to/another-script.sh`).

```shell
function main() {
  echo "Hello, world!"
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main
fi

```

For example:

```shell
#!/usr/bin/env bash
set -uo pipefail

function greet() {
  local NAME
  NAME="$1"

  if [[ -z "$NAME" ]] || [[ "$NAME" == "" ]]; then
    echo "[ERROR] No name value detected."
    return 1
  fi

  echo "Hello, $NAME!"
}

function usage() {
  echo ""
  echo "Usage: $0 [OPTIONS]"
  echo ""
  echo "Options:"
  echo "  -n, --name <name>  Name to greet"
  echo "  -h, --help         Show this help menu"
  echo ""
}

function parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -n|--name)
        NAME="$2"
        shift 2
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        echo "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
  done
}

function main() {
 parse_args "$@"

 greet "$NAME"
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi

```

### Uppercase a variable's text value

```shell
parameter=${parameter-default}; parameter=${parameter^^[a-z]}
```

## Traps

Bash [traps](https://phoenixnap.com/kb/bash-trap-command) are a feature that lets you define functionality that should occur when a specific thing happens, i.e. the script exits.

### Return to PWD on script exit

If you script uses `cd` to change your current working directory (CWD) while running, you can use this trap to return to the original path the script was called from, the $PWD:

```shell
## Save the current directory where the script was called from original_dir=$(pwd)

## Cleanup function to perform actions on script exit &
#  return to the original path
function cleanup() {
  ## You can do other things here as the script exits

  ## Return to the original path
  cd "$original_dir" || exit
}

## Set trap to run cleanup() on exit
trap cleanup EXIT

```

You can also use them more simply, like:

```shell
mkdir -p /tmp/some-path
trap rm -r /tmp/some-path EXIT

```

