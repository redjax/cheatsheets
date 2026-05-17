---
description: "[Restic](https://restic.net) is a backup/snapshotting software. It works with many backends natively, and can integrate with [rclone](https://rclone.org) to sync to many more. Tools like [resticprofile](https://creativeprojects.github.io/resticprofile/) and [Backrest](https://github.com/garethgeorge/backrest) provide scheduling & management UIs."
last_updated: "2026-05-17"
tags: ["app", "backup", "restic"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
- [Restic Examples](#restic-examples)
- [Resticprofile](#resticprofile)
  - [Resticprofile Usage](#resticprofile-usage)
  - [Resticprofile Examples](#resticprofile-examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Restic <!-- omit in toc -->

## About

[Restic](https://restic.net) is a backup/snapshotting software. It works with many backends natively, and can integrate with [rclone](https://rclone.org) to sync to many more. Tools like [resticprofile](https://creativeprojects.github.io/resticprofile/) and [Backrest](https://github.com/garethgeorge/backrest) provide scheduling & management UIs.

## Usage

| Command / Flag                      | Description                                               |
| ----------------------------------- | --------------------------------------------------------- |
| `-r <repo>`                           | Target repository (local path, sftp:, s3:, etc.).         |
| `init`                                | Initialize a new repository at the given repo path.       |
| `backup <path>`                       | Back up one or more paths into the repo.                  |
| `backup --files-from <file>`          | Read paths to back up from a file (one path per line).    |
| `backup --dry-run`                    | Show what would be backed up, without writing data.       |
| `backup --exclude <glob>`             | Exclude files matching the given glob pattern.            |
| `backup --exclude-file <file>`        | Exclude paths listed in a file (one per line).            |
| `backup --exclude-larger-than <size>` | Skip files larger than size (e.g., 1G, 100M).             |
| `backup --tag <tag>`                  | Tag snapshot with a label (can repeat for multiple tags). |
| `snapshots`                           | List all snapshots in the repo.                           |
| `snapshots --last <n>`                | Show only the last n snapshots.                           |
| `diff <id1> <id2>`                    | Show differences between two snapshot IDs.                |
| `find <pattern>`                      | Search for files matching pattern across snapshots.       |
| `restore <snapshot> --target <dir>`   | Restore a snapshot to a directory.                        |
| `restore <snapshot> --include <file>` | Restore only specific files or paths.                     |
| `mount <mountpoint>`                  | Mount the repo as a filesystem at mountpoint.             |
| `check`                               | Verify repository integrity and data consistency.         |
| `forget --keep-daily N`               | Keep last N daily snapshots, remove older ones.           |
| `forget --keep-weekly N`              | Keep last N weekly snapshots.                             |
| `forget --keep-monthly N`             | Keep last N monthly snapshots.                            |
| `forget --keep-hourly N`              | Keep last N hourly snapshots.                             |
| `forget --keep-within <duration>`     | Keep snapshots within duration (e.g., 30d).               |
| `forget --keep-tag <tag>`             | Preserve snapshots tagged with tag.                       |
| `forget --prune`                      | After forget, remove unreferenced data.                   |
| `secrets --password-file <file>`      | Read repo password from file instead of env/tty.          |

## Restic Examples

```shell
restic init -r /backup
restic backup -r /backup /home --exclude "*.log"
restic snapshots -r /backup --last 5
restic diff -r /backup abc123 def456
restic find -r /backup "*.conf"
restic restore -r /backup abc123 --target /restore
restic mount -r /backup /mnt/backup
restic check -r /backup
restic forget -r /backup --keep-daily 7 --keep-weekly 4 --prune
```

## Resticprofile

### Resticprofile Usage

| Command / Flag                 | Description                                                                         |
| ------------------------------ | ----------------------------------------------------------------------------------- |
| `-c` / `--config <path>`       | Path to `profiles.yml` configuration file.                                          |
| `--password-file <path>`         | Use a password file instead of prompting or env vars.                               |
| `--name <profile>`               | Execute related commands for a specific profile (e.g., daily-backup).               |
| `--list-profiles`                | Show all defined profiles and their default commands.                               |
| `<no command>`                   | Runs the profile’s default-command (often snapshots).                               |
| `init`                           | Initialize the repository for the selected profile (if initialize: true in config). |
| `backup`                         | Run a backup for the selected profile.                                              |
| `snapshots`                      | List snapshots in the repository (uses the profile’s repository).                   |
| `forget`                         | Run retention rules (keep-*, keep-within, prune) for the profile.                   |
| `check`                          | Run restic check on the repo for the profile.                                       |
| `mount`                          | Mount the repo according to the profile’s settings.                                 |
| `restore`                        | Restore from the repo using the profile’s repo and tags.                            |
| `schedule`                       | Install scheduler jobs (e.g., systemd timers, cron, launchd) for the profile.       |
| `schedule --all`                 | Install schedules for all profiles at once.                                         |
| `unschedule`                     | Remove scheduler jobs for the selected profile.                                     |
| `status`                         | Show status of scheduler jobs for the profile (e.g., timers, services).             |
| `run-schedule <cmd>@<profile>`   | Run a scheduled command as if invoked by the scheduler (e.g., backup@daily-backup). |
| `schedule-permission: user`      | In config, mark this profile’s schedule as user‑level instead of system‑level.      |
| `schedule-lock-wait: <duration>` | In config, max time to wait for another restic/resticprofile lock.                  |

### Resticprofile Examples

```shell
# Show all available profiles and their default commands
resticprofile --config /etc/restic/profiles.yml

# Show all profiles (even if multiple are defined)
resticprofile --config /etc/restic/profiles.yml --list-profiles

# List snapshots for the default profile
resticprofile --config /etc/restic/profiles.yml

# List snapshots for a specific profile
resticprofile --config /etc/restic/profiles.yml --name daily-backup snapshots

# Initialize a repository for a profile (if not done yet)
resticprofile --config /etc/restic/profiles.yml --name daily-backup init

# Run a backup for a given profile
resticprofile --config /etc/restic/profiles.yml --name daily-backup backup

# Run a retention / forget cleanup for a profile
resticprofile --config /etc/restic/profiles.yml --name daily-backup forget

# Install systemd / cron schedules for a specific profile
resticprofile --config /etc/restic/profiles.yml --name daily-backup schedule

# Install schedules for all profiles
resticprofile --config /etc/restic/profiles.yml schedule --all

# Check status of schedules (e.g., enabled timers/jobs)
resticprofile --config /etc/restic/profiles.yml --name daily-backup status

# Run a scheduled profile “on demand” (mimicking the schedule)
resticprofile --config /etc/restic/profiles.yml run-schedule backup@daily-backup

# Run a scheduled retention cleanup on demand
resticprofile --config /etc/restic/profiles.yml run-schedule forget@daily-backup
```

## Troubleshooting

## Links

- [Restic home](restic.net)
- [Restic docs](https://restic.readthedocs.io/en/stable/)
- [Resticprofile home](https://github.com/creativeprojects/resticprofile)
- [Resticprofile docs](https://creativeprojects.github.io/resticprofile)
- [Backrest home](https://garethgeorge.github.io/backrest/)
- [Backrest docs](https://garethgeorge.github.io/backrest/introduction/getting-started)
