---
description: "CLI controller for systemd-journal."
last_updated: "{{last_update}}"
tags: ["command", "linux", "systemd", "journalctl"]
last_updated: "2026-06-09"
---

# Journalctl <!-- omit in toc -->

## About

CLI controller for systemd-journal.

## Usage

| Command | Description |
| ------- | ----------- |
| `sudo journalctl --disk-usage` | Print how much space your systemd-journal is using. |
| `sudo journalctl --vacuum-size=500M` | Reduce logs dir to a specific size, i.e. 500 MB. |
| `sudo journalctl --vauum-time=7d` | Delete all logs older than a time period, i.e. 7 days, or 2 weeks with `--vauum-time=2weeks`, or 3 months with `--vacuum-time=6month`. |
| `sudo journalctl --vacuum-files=2` | Delete all logs from >2 boots ago. |

## Examples

## Troubleshooting

If you find yourself cleaning up systemd logs constantly, you can set some limits in `/etc/systemd/journald.conf`:

```ini
SystemMaxUse=500M
SystemKeepFree=1G
MaxRetentionSec=3month
```

After changing this file, restart the systemd-journal service with `sudo systemctl restart systemd-journald`.

## Links
