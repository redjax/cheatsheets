---
description: "Rsync is a fast remote and local file-copying tool."
last_updated: "{{last_update}}"
tags: ["command", ]
---
## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Rsync <!-- omit in toc -->

[rsync](https://linux.die.net/man/1/rsync) is a command line utility for copying files quickly and efficiently, both locally and between remotes.

> [!NOTE]
> `rsync` must be installed on the remote machine to use the command to copy files to/from a remote.

## Usage

General command form (copy a file/path to a remote):

```shell
rsync -azvh --progress /local/path user@remote:/remote/path
```

Args:

| arg | description |
| --- | ----------- |
| `-r` | Recursive copy (unnecessary with `-a`) |
| `-a` | Archive mode, includes recursive transfer |
| `-z` | Compress the data |
| `-v` | Verbose/detailed info during transfer |
| `-h` | Human-readable output |

## Examples

## Troubleshooting

## Links
