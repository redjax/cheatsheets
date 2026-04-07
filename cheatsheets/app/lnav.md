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

## Examples

## Troubleshooting

## Links

- [lnav home](https://lnav.org)
- [lnav Github](https://github.com/tstack/lnav)
- [lnav docs](https://docs.lnav.org/en/v0.13.1/)

