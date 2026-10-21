---
description: "[Gitu](https://github.com/altsem/gitu) is a TUI for Git porcelain operations, inspired by Magit for Emacs."
last_updated: "2026-10-21"
tags: ["app", "git", "tui"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Gitu <!-- omit in toc -->

## About

[Gitu](https://github.com/altsem/gitu) is a TUI for Git porcelain operations, inspired by Magit for Emacs.

## Usage

| **Context**      | Key(s)                    | Action / option                            |
| ---------------- | ------------------------- | ------------------------------------------ |
| Root             | `q`, `esc`                | Quit                                       |
| Root             | `g+r`                     | Refresh                                    |
| Root             | `tab`                     | Toggle section                             |
| Root             | `k`, `up`                 | Move up                                    |
| Root             | `j`, `down`               | Move down                                  |
| Root             | `ctrl+k`, `ctrl+up`       | Move up one line                           |
| Root             | `ctrl+j`, `ctrl+down`     | Move down one line                         |
| Root             | `alt+k`, `alt+up`         | Move to previous section                   |
| Root             | `alt+j`, `alt+down`       | Move to next section                       |
| Root             | `alt+h`, `alt+left`       | Move to parent section                     |
| Root             | `ctrl+u`                  | Move up half a page                        |
| Root             | `ctrl+d`                  | Move down half a page                      |
| Root             | `g+g`                     | Move to top                                |
| Root             | `G`                       | Move to bottom                             |
| Root             | `ctrl+y`                  | Scroll view up                             |
| Root             | `ctrl+e`                  | Scroll view down                           |
| Root             | `/`                       | Search                                     |
| Root             | Unbound                   | Search backward                            |
| Root             | `n`                       | Next search result                         |
| Root             | `N`                       | Previous search result                     |
| Root             | `Y`                       | Show refs                                  |
| Root             | `enter`                   | Show selection                             |
| Root             | `K`                       | Discard                                    |
| Root             | `s`                       | Stage                                      |
| Root             | `u`                       | Unstage                                    |
| Root             | `a`                       | Apply                                      |
| Root             | `v`                       | Reverse                                    |
| Root             | `y`                       | Copy hash                                  |
| Root             | `B`                       | Blame                                      |
| Root             | `h`, `?`                  | Open help menu                             |
| Root             | `b`                       | Open branch menu                           |
| Root             | `c`                       | Open commit menu                           |
| Root             | `f`                       | Open fetch menu                            |
| Root             | `l`                       | Open log menu                              |
| Root             | `m`                       | Open merge menu                            |
| Root             | `F`                       | Open pull menu                             |
| Root             | `P`                       | Open push menu                             |
| Root             | `r`                       | Open rebase menu                           |
| Root             | `M`                       | Open remote menu                           |
| Root             | `X`                       | Open reset menu                            |
| Root             | `V`                       | Open revert menu                           |
| Root             | `A`                       | Open cherry-pick menu                      |
| Root             | `z`                       | Open stash menu                            |
| Picker           | `down`, `ctrl+n`, `tab`   | Next item                                  |
| Picker           | `up`, `ctrl+p`, `backtab` | Previous item                              |
| Picker           | `enter`                   | Confirm                                    |
| Picker           | `esc`, `ctrl+c`           | Cancel                                     |
| Help menu        | `q`, `h`, `?`, `esc`      | Close menu                                 |
| Branch menu      | `b`                       | Checkout                                   |
| Branch menu      | `c`                       | Checkout new branch                        |
| Branch menu      | `s`                       | Spin off branch                            |
| Branch menu      | `m`                       | Rename branch                              |
| Branch menu      | `K`                       | Delete branch                              |
| Branch menu      | `q`, `esc`                | Close menu                                 |
| Commit menu      | `-a`                      | Option: `--all`                            |
| Commit menu      | `-e`                      | Option: `--allow-empty`                    |
| Commit menu      | `-v`                      | Option: `--verbose`                        |
| Commit menu      | `-n`                      | Option: `--no-verify`                      |
| Commit menu      | `-R`                      | Option: `--reset-author`                   |
| Commit menu      | `-s`                      | Option: `--signoff`                        |
| Commit menu      | `c`                       | Commit                                     |
| Commit menu      | `a`                       | Amend commit                               |
| Commit menu      | `e`                       | Extend commit                              |
| Commit menu      | `f`                       | Fixup commit                               |
| Commit menu      | `F`                       | Instant fixup commit                       |
| Commit menu      | `q`, `esc`                | Close menu                                 |
| Fetch menu       | `-p`                      | Option: `--prune`                          |
| Fetch menu       | `-t`                      | Option: `--tags`                           |
| Fetch menu       | `a`                       | Fetch all                                  |
| Fetch menu       | `e`                       | Fetch from elsewhere                       |
| Fetch menu       | `p`                       | Fetch from push remote                     |
| Fetch menu       | `u`                       | Fetch from upstream                        |
| Fetch menu       | `q`, `esc`                | Close menu                                 |
| Log menu         | `l`                       | Log current                                |
| Log menu         | `o`                       | Log other                                  |
| Log menu         | `-n`                      | Option: `-n`                               |
| Log menu         | `-F`                      | Option: `--grep`                           |
| Log menu         | `q`, `esc`                | Close menu                                 |
| Merge menu       | `-f`                      | Option: `--ff-only`                        |
| Merge menu       | `-n`                      | Option: `--no-ff`                          |
| Merge menu       | `a`                       | Abort merge                                |
| Merge menu       | `c`                       | Continue merge                             |
| Merge menu       | `m`                       | Merge                                      |
| Merge menu       | `q`, `<esc>`              | Close menu (`<esc>` as supplied in config) |
| Pull menu        | `-r`                      | Option: `--rebase`                         |
| Pull menu        | `p`                       | Pull from push remote                      |
| Pull menu        | `u`                       | Pull from upstream                         |
| Pull menu        | `e`                       | Pull from elsewhere                        |
| Pull menu        | `q`, `esc`                | Close menu                                 |
| Push menu        | `-f`                      | Option: `--force-with-lease`               |
| Push menu        | `-F`                      | Option: `--force`                          |
| Push menu        | `-h`                      | Option: `--no-verify`                      |
| Push menu        | `-n`                      | Option: `--dry-run`                        |
| Push menu        | `p`                       | Push to push remote                        |
| Push menu        | `u`                       | Push to upstream                           |
| Push menu        | `e`                       | Push to elsewhere                          |
| Push menu        | `q`, `esc`                | Close menu                                 |
| Rebase menu      | `-k`                      | Option: `--keep-empty`                     |
| Rebase menu      | `-p`                      | Option: `--preserve-merges`                |
| Rebase menu      | `-d`                      | Option: `--committer-date-is-author-date`  |
| Rebase menu      | `-a`                      | Option: `--autosquash`                     |
| Rebase menu      | `-A`                      | Option: `--autostash`                      |
| Rebase menu      | `-i`                      | Option: `--interactive`                    |
| Rebase menu      | `-h`                      | Option: `--no-verify`                      |
| Rebase menu      | `i`                       | Interactive rebase                         |
| Rebase menu      | `a`                       | Abort rebase                               |
| Rebase menu      | `c`                       | Continue rebase                            |
| Rebase menu      | `e`                       | Rebase elsewhere                           |
| Rebase menu      | `f`                       | Autosquash rebase                          |
| Rebase menu      | `q`, `esc`                | Close menu                                 |
| Remote menu      | `a`                       | Add remote                                 |
| Remote menu      | `K`                       | Remove remote                              |
| Remote menu      | `r`                       | Rename remote                              |
| Remote menu      | `q`, `esc`                | Close menu                                 |
| Reset menu       | `s`                       | Soft reset                                 |
| Reset menu       | `m`                       | Mixed reset                                |
| Reset menu       | `h`                       | Hard reset                                 |
| Reset menu       | `q`, `esc`                | Close menu                                 |
| Revert menu      | `-e`                      | Option: `--edit`                           |
| Revert menu      | `-E`                      | Option: `--no-edit`                        |
| Revert menu      | `-s`                      | Option: `--signoff`                        |
| Revert menu      | `a`                       | Abort revert                               |
| Revert menu      | `c`                       | Continue revert                            |
| Revert menu      | `V`                       | Revert commit                              |
| Revert menu      | `q`, `esc`                | Close menu                                 |
| Cherry-pick menu | `-n`                      | Option: `--no-commit`                      |
| Cherry-pick menu | `-s`                      | Option: `--signoff`                        |
| Cherry-pick menu | `-e`                      | Option: `--edit`                           |
| Cherry-pick menu | `a`                       | Abort cherry-pick                          |
| Cherry-pick menu | `c`                       | Continue cherry-pick                       |
| Cherry-pick menu | `A`                       | Cherry-pick                                |
| Cherry-pick menu | `q`, `esc`                | Close menu                                 |
| Stash menu       | `-a`                      | Option: `--all`                            |
| Stash menu       | `-u`                      | Option: `--include-untracked`              |
| Stash menu       | `z`                       | Stash                                      |
| Stash menu       | `i`                       | Stash index                                |
| Stash menu       | `w`                       | Stash worktree                             |
| Stash menu       | `x`                       | Stash, keeping index                       |
| Stash menu       | `p`                       | Pop stash                                  |
| Stash menu       | `a`                       | Apply stash                                |
| Stash menu       | `k`                       | Drop stash                                 |
| Stash menu       | `q`, `esc`                | Close menu                                 |

## Examples

## Troubleshooting

## Links
