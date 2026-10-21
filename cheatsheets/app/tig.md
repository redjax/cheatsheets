---
description: "[tig](https://github.com/jonas/tig) is a text-mode CLI for git. Unlike similar tools like LazyGit and gitu, tig is an ncurses-based interface that functions mainly as a repository browsers. Its commands wrap the `git` CLI to simplify and enhance Git operations, particulary around `log` and `reflog`."
last_updated: "2026-10-21"
tags: ["app", "git", "cli"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Tig <!-- omit in toc -->

## About

[tig](https://github.com/jonas/tig) is a text-mode CLI for git. Unlike similar tools like LazyGit and gitu, tig is an ncurses-based interface that functions mainly as a repository browsers. Its commands wrap the `git` CLI to simplify and enhance Git operations, particulary around `log` and `reflog`.

📖 [`tig` Manual](https://jonas.github.io/tig/doc/manual.html)

## Usage

| **Context**     | Key(s)                 | Action                                                                                                 |
| --------------- | ---------------------- | ------------------------------------------------------------------------------------------------------ |
| View switching  | `m`                    | Open main view (commit summaries)                                                                      |
| View switching  | `d`                    | Open diff view                                                                                         |
| View switching  | `l`                    | Open log view                                                                                          |
| View switching  | `p`                    | Open pager view                                                                                        |
| View switching  | `t`                    | Open tree view (directories and files)                                                                 |
| View switching  | `f`                    | Open blob view (file contents)                                                                         |
| View switching  | `g`                    | Open grep view                                                                                         |
| View switching  | `b`                    | Open blame view                                                                                        |
| View switching  | `r`                    | Open refs view (branches, remotes, and tags)                                                           |
| View switching  | `y`                    | Open stash view                                                                                        |
| View switching  | `h`                    | Open help / keybindings                                                                                |
| View switching  | `s`                    | Open status view                                                                                       |
| View switching  | `c`                    | Open stage view                                                                                        |
| View management | `q`                    | Close current view; quit if it is the last view                                                        |
| View management | `Q`                    | Quit all views                                                                                         |
| View management | `tab`                  | Switch focus to next view                                                                              |
| View management | `R`                    | Reload / refresh current view                                                                          |
| View management | `O`                    | Maximize current view                                                                                  |
| Main / log      | `enter`                | Open selected commit’s diff in a split view                                                            |
| Diff            | `enter`                | Scroll down one line                                                                                   |
| Navigation      | `k`                    | Move cursor up one line                                                                                |
| Navigation      | `j`                    | Move cursor down one line                                                                              |
| Navigation      | `up`, `down`           | Move up / down; in a split view, may select the previous / next parent-view commit and update its diff |
| Navigation      | `pageup`, `-`          | Move up one page                                                                                       |
| Navigation      | `pagedown`, `space`    | Move down one page                                                                                     |
| Navigation      | `home`                 | Jump to first line                                                                                     |
| Navigation      | `end`                  | Jump to last line                                                                                      |
| Navigation      | `H`                    | Go to `HEAD` commit                                                                                    |
| Tree            | `,`                    | Go to parent directory                                                                                 |
| Blame           | `,`                    | Load blame for parent commit; prompt for parent on merges                                              |
| Scrolling       | `insert`               | Scroll up one line                                                                                     |
| Scrolling       | `delete`               | Scroll down one line                                                                                   |
| Scrolling       | `ScrollBack`           | Scroll up one page                                                                                     |
| Scrolling       | `ScrollFwd`            | Scroll down one page                                                                                   |
| Scrolling       | `left`                 | Scroll left                                                                                            |
| Scrolling       | `right`                | Scroll right                                                                                           |
| Scrolling       | `\|`                   | Scroll to first column                                                                                 |
| Search          | `/`                    | Search forward using a regular expression                                                              |
| Search          | `?`                    | Search backward using a regular expression                                                             |
| Search          | `n`                    | Next search match                                                                                      |
| Search          | `N`                    | Previous search match                                                                                  |
| Status          | `u`                    | Stage / unstage selected file, or add an untracked file                                                |
| Stage           | `u`                    | Stage / unstage current diff chunk when on its header; otherwise update the entire displayed diff      |
| Stage           | `1`                    | Stage a single diff line                                                                               |
| Stage           | `@`                    | Jump to next diff chunk                                                                                |
| Status          | `M`                    | Resolve selected unmerged file with `git mergetool`                                                    |
| Status / stage  | `!`                    | Discard unstaged file changes (destructive)                                                            |
| Status          | `C`                    | Run `git commit`                                                                                       |
| Main            | `C`                    | Cherry-pick selected commit                                                                            |
| Diff / stage    | `]`                    | Increase diff context                                                                                  |
| Diff / stage    | `[`                    | Decrease diff context                                                                                  |
| Display         | `ctrl+l`               | Redraw screen                                                                                          |
| Display         | `#`                    | Toggle line numbers                                                                                    |
| Display         | `D`                    | Cycle date display modes                                                                               |
| Display         | `A`                    | Cycle author display modes                                                                             |
| Display         | `G`                    | Toggle revision graph                                                                                  |
| Display         | `~`                    | Cycle line graphics mode                                                                               |
| Display         | `F`                    | Toggle reference names (branches and tags)                                                             |
| Display         | `W`                    | Toggle ignoring whitespace in diffs                                                                    |
| Display         | `X`                    | Toggle commit ID display                                                                               |
| Display         | `%`                    | Toggle file filtering: selected file’s diff versus full diff                                           |
| Display         | `$`                    | Toggle highlighting of commit-title overflow                                                           |
| Miscellaneous   | `o`                    | Open option menu                                                                                       |
| Miscellaneous   | `e`                    | Open file in editor                                                                                    |
| Miscellaneous   | `z`                    | Stop background loading                                                                                |
| Miscellaneous   | `v`                    | Show version                                                                                           |
| Prompt          | `:`                    | Open command prompt                                                                                    |
| Prompt          | `:<number>`            | Jump to line number, e.g. `:80`                                                                        |
| Prompt          | `:<sha>`               | Jump to commit, e.g. `:2f12bcc`                                                                        |
| Prompt          | `:<key>`               | Execute a keybinding, e.g. `:q`                                                                        |
| Prompt          | `:<action>`            | Execute a Tig action, e.g. `:edit`                                                                     |
| Prompt          | `:!<command>`          | Run a system command and show output in a pager, e.g. `:!git log -p`                                   |
| Prompt          | `:goto <rev>`          | Jump to revision, e.g. `:goto main`                                                                    |
| Prompt          | `:save-display <file>` | Save current display to a file                                                                         |
| Prompt          | `:save-options <file>` | Save current options to a file                                                                         |
| Prompt          | `:script <file>`       | Execute Tig commands from a file                                                                       |

## Examples

## Troubleshooting

## Links
