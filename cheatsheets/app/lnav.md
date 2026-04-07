---
description: "[lnav](https://lnav.org/) is a logfile navigator/TUI. It is extremely versatile at parsing logs from different sources/filetypes."
last_updated: "2026-04-07"
tags: ["app", "cli", "tui", "logging", "linux"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
  - [Keybinds Cheatsheet](#keybinds-cheatsheet)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# LNav - The Logfile Navigator <!-- omit in toc -->

## About

[lnav](https://lnav.org/) is a logfile navigator/TUI. It is extremely versatile at parsing logs from different sources/filetypes.

## Usage

### Keybinds Cheatsheet

| Key | Description |
| --- | ----------- |
| `e` | Move to next error |
| `Shift+e` | Move to previous error |
| `w` | Move to next warning |
| `Shift+w` | Move to previous warning |
| `/` | Start a text search in the current log. Use `TAB` for completion from text in session |
| `0-9` | Move messages on screen in 10 minute intervals. e.g. `2` moves to the first message after the next 20 minute mark, `3` moves to the next half hour mark, etc. |
| `` ` `` | Focus the breadcrumb bar. Use `TAB` or right arrow to move between crumbs |
| `:` | Activate command prompt for commands like `:goto` |
| `:goto` | Go to a given timestamp/line number |
| `?`      | Open built-in help                                                   |
| `q`      | Go back / quit the current view                                      |
| `Q`      | Go back while keeping times aligned between views                    |
| `Ctrl+L` | Refresh/redraw the screen                                            |
| `SPACE`  | Page down                                                            |
| `b`      | Page up                                                              |
| `n`      | Next search match                                                    |
| `N`      | Previous search match                                                |
| `m`      | Mark or unmark the current line                                      |
| `M`      | Mark or unmark a range of lines                                      |
| `c`      | Copy marked lines to the clipboard                                   |
| `C`      | Clear marked lines                                                   |
| `g`      | Jump to the top of the current view                                  |
| `G`      | Jump to the bottom of the current view                               |
| `[ / ]`  | Move through time buckets or nearby boundaries, depending on context |
| `TAB`    | Cycle through completions in prompts and search fields               |
| `:`      | Open the command prompt                                              |
| `/`      | Start a regex/text search                                            |
| `;`      | Enter SQL query mode                                                 |
| `:goto 2026-04-07T13:00:00`      | Jump to a specific timestamp.                                              |
| `:goto 12345`                    | Jump to a line number.                                                     |
| `:pipe-to grep ERROR`            | Pipe bookmarked lines to a shell command and open the result in lnav lnav. |
| `:pipe-line-to awk '{print $1}'` | Pipe the current line to a shell command and open the result in lnav lnav. |
| `:write-csv-to -`                | Export query results as CSV to stdout lnav.                                |
| `:write-json-to results.json`    | Export query results as JSON lnav.                                         |
| `:redirect-to output.txt`        | Send command output to a file lnav.                                        |

## Examples

| Command                                            | Description                                                                                                                 |
| -------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `lnav /var/log/syslog`                               | Open a single log file.                                                                                                     |
| `lnav app.log error.log`                             | Open multiple files and merge them into one timeline view. lnav will interleave entries by timestamp when possible mankier. |
| `lnav /var/log`                                      | Open every supported file in a directory and keep watching for new log files lnav.                                          |
| `lnav ~/logs/*.log`                                  | Open a set of matching files with a shell glob.                                                                             |
| `journalctl -u nginx -f \\\| lnav`                   | Pipe command output into lnav for interactive viewing.                                                                      |
| `docker compose logs -f \\\| lnav`                   | Use lnav to inspect live container logs.                                                                                    |
| `kubectl logs deploy/api -f \\\| lnav`               | View Kubernetes logs in lnav.                                                                                               |
| `tail -f /var/log/auth.log \\\| lnav`                | Follow a live log stream in lnav.                                                                                           |
| `ssh user@host 'journalctl -u nginx -n 200' \\\| lnav` | Read remote logs through SSH.                                                                                               |
| `lnav -C /var/log`                                   | Validate log formats in a directory without opening the full interactive UI mankier.                                        |

## Troubleshooting

## Links

- [lnav home](https://lnav.org)
- [lnav Github](https://github.com/tstack/lnav)
- [lnav docs](https://docs.lnav.org/en/v0.13.1/)

