---
description: "[Rclone](https://rclone.org) is a CLI utility to manage files on cloud storage. It can interact with [over 70 cloud storage products](https://rclone.org/#providers), and can act as a bridge for tools like Restic."
last_updated: "2026-05-27"
tags: ["app", "backup", "sync", "cli"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Setup](#setup)
  - [Install Rclone](#install-rclone)
- [Usage](#usage)
- [Examples](#examples)
  - [Bandwidth limit examples](#bandwidth-limit-examples)
- [Troubleshooting](#troubleshooting)
- [Link](#link)

# Rclone <!-- omit in toc -->

## About

[Rclone](rclone.org) is a CLI utility to manage files on cloud storage. It can interact with [over 70 cloud storage products](https://rclone.org/#providers), and can act as a bridge for tools like Restic.

## Setup

### Install Rclone

```shell
curl https://rclone.org/install.sh | sudo bash
```

## Usage

| Command / Arg                          | Description                                                           |
| -------------------------------------- | --------------------------------------------------------------------- |
| `config`                               | Start interactive setup for remotes                                   |
| `listremotes`                          | Show all configured remotes                                           |
| `lsd <remote>:`                        | List top-level directories in a remote                                |
| `ls <remote:path>`                     | List files with sizes                                                 |
| `lsf <remote:path>`                    | Fast listing of files (good for scripts / connectivity check)         |
| `copy src dest`                        | Copy files (no deletions)                                             |
| `sync src dest`                        | Mirror source to destination (deletes extras)                         |
| `move src dest`                        | Move files (deletes from source after copy)                           |
| `check src dest`                       | Compare files between locations (no changes)                          |
| `check src dest --one-way`             | Only verify src → dest consistency (backup validation)                |
| `check src dest --size-only`           | Compare only file sizes (ignore hashes/mod times)                     |
| `check src dest --download`            | Verify by downloading destination data (strongest validation, slower) |
| `check src dest --differ`              | Show only files that differ                                           |
| `check src dest --missing-on-dst`      | Show files present in source but missing on destination               |
| `check src dest --missing-on-src`      | Show files present in destination but missing in source               |
| `check src dest --combined`            | Show a combined report of all differences                             |
| `check src dest --one-way --size-only` | Lightweight backup check (fast but weaker guarantee)                  |
| `mkdir <remote:path>`                  | Create directory in remote                                            |
| `delete <remote:path>`                 | Delete files in a path                                                |
| `purge <remote:path>`                  | Delete entire directory tree                                          |
| `[-v \| -vv]`                          | Increase logging verbosity (debug-level with `-vv`)                   |
| `config reconnect <remote>:`           | Reauthenticate a remote                                               |
| `--dry-run`                            | Simulate operation without making changes                             |
| `--verbose`                            | Show detailed output (same family as `-v`)                            |
| `--progress`                           | Show transfer progress                                                |
| `--bwlimit <rate>`                     | Limit bandwidth usage. Supports formats like `10M`, `1M:08:00,10M:17:00` (different limits per time), `on`/`off`, or burst-style limits like `1M:off` |
| `--transfers N`                        | Number of concurrent file transfers                                   |
| `--delete-during`                      | Delete files during sync (faster, less space used)                    |
| `--log-file <file>`                    | Write logs to a file instead of stdout                                |
| `--stats 5s`                           | Print periodic transfer stats every interval                          |
| `ncdu <remote:path>` | Interactive disk usage explorer |

## Examples

### General usage examples

`rclone` syntax:

```shell
rclone <command> [flags] <source> <destination>
```

Common patterns:

```shell
rclone copy   <src> <remote:dest>
rclone sync   <src> <remote:dest>
rclone move   <src> <remote:dest>

rclone ls     <remote:path>
rclone lsd    <remote:>
rclone lsf    <remote:path>
```

Source and destination formats:

| Type          | Example                     |
| ------------- | --------------------------- |
| Local path    | `~/data`, `/home/user/data` |
| Remote root   | `remote:`                   |
| Remote folder | `remote:backups`            |
| Nested path   | `remote:backups/2026`       |

### Common operations

View/edit `rclone`'s config:

```shell
rclone config
```

Show configured remotes:

```shell
rclone listremotes
```

List the top-level directories in a remote:

```shell
rclone lsd remote-name:
```

List top-level directories in a remote starting at a given path:

```shell
rclone ls remote-name:backup/path
```

### Rclone copy examples

| Example                                               | Meaning                                          |
| ----------------------------------------------------- | ------------------------------------------------ |
| `rclone copy ~/data remote:backup`                    | Basic upload (no deletions on destination)       |
| `rclone copy ~/data remote:backup --progress`         | Show transfer progress                           |
| `rclone copy ~/data remote:backup --dry-run`          | Simulate what would be copied                    |
| `rclone copy ~/data remote:backup --ignore-existing`  | Skip files that already exist on destination     |
| `rclone copy ~/data remote:backup --update`           | Only copy newer files (based on mod time)        |
| `rclone copy ~/data remote:backup --exclude "*.tmp"`  | Exclude temp files                               |
| `rclone copy ~/data remote:backup --filter "- *.log"` | Filter rule version of exclusion                 |
| `rclone copy ~/data remote:backup --checksum`         | Compare using checksums instead of mod time/size |


### Rclone sync examples

| Example                                                               | Meaning                                                               |
| --------------------------------------------------------------------- | --------------------------------------------------------------------- |
| `rclone sync ~/data remote:backup`                                    | Mirror local -> remote (deletes extras on destination)                 |
| `rclone sync ~/data remote:backup --dry-run`                          | Preview destructive changes safely                                    |
| `rclone sync ~/data remote:backup --progress`                         | Show progress during sync                                             |
| `rclone sync ~/data remote:backup --delete-excluded`                  | Delete files that match exclude rules                                 |
| `rclone sync ~/data remote:backup --exclude "*.tmp"`                  | Exclude files and remove excluded ones from destination               |
| `rclone sync ~/data remote:backup --backup-dir remote:backup-archive` | Move overwritten/deleted files into backup folder instead of deleting |
| `rclone sync ~/data remote:backup --update`                           | Only overwrite older files on destination                             |
| `rclone sync ~/data remote:backup --size-only`                        | Compare only file sizes (faster, weaker consistency)                  |
| `rclone sync ~/data remote:backup --checksum`                         | Strong verification using checksums                                   |


### Bandwidth limit examples

| Example | Meaning |
|--------|--------|
| `--bwlimit 10M` | Limit to 10 MB/s always |
| `--bwlimit 1M:08:00,10M:17:00` | 1 MB/s during work hours, 10 MB/s after 5pm |
| `--bwlimit 1M:off` | 1 MB/s during transfers, unlimited at other times |
| `--bwlimit off` | No bandwidth limit |

## Troubleshooting

## Link

- [Rclone home](https://rclone.org)
- [Rclone docs](https://rclone.org/docs/)
  - [Rclone commands reference](https://rclone.org/commands/)
  - [Rclone providers reference](https://rclone.org/#providers)
