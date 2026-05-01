---
description: "Linux OS."
last_updated: "2026-05-01"
tags: ["linux", "system", "os"]
---

# Linux <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

## Paths

### Variables

| Variable                                                  | Description                             |
| --------------------------------------------------------- | --------------------------------------- |
| `THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" | Get parent directory of current script. |
| `REL_PATH=$(realpath -m "${THIS_DIR}/../..")              | Get absolute path from a relative path. |

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

