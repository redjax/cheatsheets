---
description: "[tmux](https://github.com/tmux/tmux) is a terminal multiplexer, allowing you to split 1 terminal session into many with windows and panes."
last_updated: "2026-04-07"
tags: ["app", "cli", "tui", "linux"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Tmux - Terminal Multiplexer <!-- omit in toc -->

## About

[tmux](https://github.com/tmux/tmux) is a terminal multiplexer, allowing you to split 1 terminal session into many with windows and panes.

## Usage

### Keybinds Cheatsheet

> [!NOTE]
> `C-b` is the default prefix for Tmux; wherever you see `C-*`, that means "hold CONTROL and press a key." In this case, `C-b` means `CTRL+b`.
>
> If you rebind the default `C-b` prefix, make sure to use your prefix in any of the commands below. For example, if you set the prefix to
> `C-a`, wherever you see `C-b` in the table below, use `C-a` instead.

| Key        | Description                                                           |
| ---------- | --------------------------------------------------------------------- |
| `C-b`     | Default prefix key. Press this first before most tmux actions. |
| `C-b ?`   | Show all key bindings.                                         |
| `C-b c`   | Create a new window.                                             |
| `C-b w`   | Choose a window interactively.                                   |
| `C-b n`   | Move to the next window.                                         |
| `C-b p`   | Move to the previous window.                                     |
| `C-b 0-9` | Switch to window 0 through 9.                                    |
| `C-b ,`   | Rename the current window.                                       |
| `C-b $`   | Rename the current session.                                      |
| `C-b d`   | Detach from the current session.                                 |
| `C-b s`   | Choose a session interactively.                                  |
| `C-b [`   | Enter copy mode to scroll/search/select text.                    |
| `C-b ]`   | Paste the most recently copied buffer.                           |
| `C-b %`   | Split the current pane vertically.                               |
| `C-b "`   | Split the current pane horizontally.                             |
| `C-b o`   | Move to the next pane.                                           |
| `C-b ;`   | Jump back to the previously active pane.                         |
| `C-b x`   | Kill the current pane.                                           |
| `C-b z`   | Toggle zoom for the current pane.                                |
| `C-b !`   | Break the current pane into its own window.                      |
| `C-b q`   | Briefly display pane numbers.                                    |
| `C-b t`   | Show the time.                                                   |
| `C-b :`   | Open the command prompt.                                         |
| `C-b ?`   | Show the key binding list.                                       |
| `C-b {`   | Swap the current pane with the previous pane.                    |
| `C-b }`   | Swap the current pane with the next pane.                        |
| `C-b ,`   | Rename window, handy for tagging work like logs, api, or deploy. |
| `C-b L`   | Switch back to the last session.                                 |
| `C-b m`   | Mark the current pane.                                           |
| `C-b M`   | Clear the marked pane.                                           |
| `C-b C-o` | Rotate panes in the current window.                              |
| `C-b C-z` | Suspend the client.                                              |

### Commands Cheatsheet

| Command | Description |
| `tmux ls` | List sessions |
| `tmux kill-server` | Kill all sessions |
| `tmux a -t <index-or-name>` | Attach to a session by index or name |

## Examples

Start a session:

```shell
tmux
tmux new
```

Start a new named session:

```shell
tmux new -s session-name
```

Attach to a session:

```shell
tmux attach
tmux attach -t session-name
tmux a -t session-name
```

Manage sessions:

```shell
tmux ls # list available sessions
tmux kill-session -t session-name
tmux switch -t session2-name
```

Work with windows:

```shell
tmux new-window
tmux new-window -n window-name
tmux rename-window window-name
tmux select-window -t 2 # select a window by index/name
```

Work with panes:

```shell
tmux split-window
tmux split-window -h
tmux split-window -v
tmux select-pane -L
tmux select-pane -R
tmux select-pane -U
tmux select-pane -D
tmux resize-pane -L 10
```

## Troubleshooting

## Links
