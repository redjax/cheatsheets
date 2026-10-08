---
description: ""
last_updated: "{{last_update}}"
tags: ["command", ]
last_updated: "2026-10-08"
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Atd <!-- omit in toc -->

## About

Atd is a job scheduler daemon that runs jobs scheduled for later execution. Like cron, but when you only need to run a script on a schedule a single time.

## Installation

- Debian/Ubuntu:

  ```shell
  sudo apt update -y
  sudo apt install -y at
  sudo systemctl enable --now atd
  ```

- RedHat/Fedora/CentOS/AlmaLinux/RockyLinux:

  ```shell
  # DNF
  sudo dnf install -y at
  sudo systemctl enable --now atd

  # YUM
  sudo yum install at
  sudo systemctl enable --now atd
  ```

- openSUSE:

  ```shell
  sudo zypper install at
  sudo systemctl enable --now atd
  ```

- Arch Linux:

  ```shell
  sudo pacman -S at
  sudo systemctl enable --now atd
  ```

## Usage

| Task | Command / Example | Notes |
| --- | --- | --- |
| Install - Debian/Ubuntu | `sudo apt install at` | Installs `at` and `atd` |
| Install - RHEL/Fedora/CentOS/Alma/Rocky | `sudo dnf install at` | Use `yum` on older systems |
| Install - openSUSE | `sudo zypper install at` |  |
| Install - Arch | `sudo pacman -S at` |  |
| Start`atd` | `sudo systemctl start atd` | Start immediately |
| Enable`atd`at boot | `sudo systemctl enable atd` | Start automatically after reboot |
| Enable + start | `sudo systemctl enable --now atd` | Usually the easiest option |
| Check`atd`status | `systemctl status atd` |  |
| Restart`atd` | `sudo systemctl restart atd` |  |
| View`atd`logs | `journalctl -u atd` |  |
| Follow logs | `journalctl -u atd -f` | Live log output |
| Schedule a command | `echo 'command' \| at 3:00 AM` | Runs once |
| Schedule Bash command | `echo '/usr/bin/bash -lc "command"' \| at 3:00 AM` | Explicitly use Bash |
| Schedule a script | `echo '/usr/bin/bash /path/script.sh' \| at 3:00 AM` |  |
| Multiple commands | `at 3:00 AM <<'EOF'`\<br\>`command1`\<br\>`command2`\<br\>`EOF` | Heredoc is cleaner for multiple commands |
| Run in a specific directory | `echo 'cd /path && command' \| at 3:00 AM` | Don't assume the working directory |
| Schedule relative time | `echo 'command' \| at now + 30 minutes` |  |
| Schedule tomorrow | `echo 'command' \| at 3:00 AM tomorrow` |  |
| Schedule as root | `sudo at 3:00 AM` | Job runs as root; no sudo prompt at execution time |
| List your jobs | `atq` | Shows job IDs and scheduled times |
| List root's jobs | `sudo atq` |  |
| Cancel a job | `atrm JOB_ID` | Example: `atrm 12` |
| Cancel root job | `sudo atrm JOB_ID` |  |
| Inspect a job | `at -c JOB_ID` | Shows what will actually execute |
| Redirect output | `echo 'command >> /tmp/job.log 2>&1' \| at 3:00 AM` | Useful for troubleshooting |
| Stop after a failure | `set -e` | Useful inside a Bash heredoc/script |
| Check command path | `which command` | `at` has a more limited environment than your interactive shell |
| Check permissions | `/etc/at.allow` / `/etc/at.deny` | Controls who can use `at` |
| Manual | `man at` | Full `at` documentation |
| Daemon manual | `man atd` | `atd` daemon documentation |

## Examples

## Troubleshooting

## Links
