---
description: "Grep is used to find text in a file or files within a directory path."
last_updated: "{{last_update}}"
tags: ["command", "cli"]
last_updated: "2026-05-06"
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Grep <!-- omit in toc -->

## About

Grep is used to find text in a file or files within a directory path.

## Usage

General syntax:

```shell
grep [OPTIONS] "search_string" [path/to/search]
```

| Command | Description |
| ------- | ----------- |
| `grep -i [string]` | Case-insensitive search |
| `grep -w [string]` | Search full words |
| `grep -A [num]`    | Show N lines before match |
| `grep -B [num]`    | Show N lines after match |
| `grep -C [num]`   | Show N lines around match |
| `grep -r [string] /path/to/search` | Recursively search files in a path |
| `grep -R [string] /path/to/search` | Follow symlinks |
| `grep -v [string]` | Return all lines which DON'T match a pattern |
| `grep -e "^regex-string"` | Use regex |
| `grep -E 're(gex|GEX)' string` | Extended regex |
| `grep -c [string]` | Count number of matches |
| `grep -l [string]` | Print the name of the file(s) of matches |
| `grep -L [string]` | Show non-matching filenames |
| `grep -o [string]` | Only show matching part of the string |
| `grep -n [string]` | Show line number for matches |
| `grep --include="*.fileext"` | Only search files with specific filetype |
| `grep --exclude="*.fileext"` | Do not search files with specific filetype |
| `grep -m [num] [string]` | Stop search after N matches |

## Examples

## Troubleshooting

## Links

