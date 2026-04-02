---
description: "Terminal-based IDE, next-gen Vi(m)."
last_updated: "2026-04-02"
tags: ["neovim", "app", "cli", "tui"]
---

# Neovim <!-- omit in toc -->

[https://neovim.io/](https://neovim.io/)

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
  - [CLI args](#cli-args)
  - [Navigation](#navigation)
  - [Search](#search)

## Usage

### CLI args

Launching Neovim with CLI args controls how the program opens. Here are some useful quick-commands you can use:

- Open Neovim to a specific line

  ```shell
  neovim +<number>
  ```

  - Example: Open Neovim to line 113: `neovim +113`
- Run Lazy package manager sync
  ```shell
  nvim --headless "+Lazy! sync" +qa
  ```

### Navigation

When editing a file, these less obvious keybinds are useful for navigating around the file/lines:

| Keybind | Description |
| ------- | ----------- |
| `CTRL+e` | scroll page down without scrolling cursor |
| `CTRL+y` | scroll page up without scrolling cursor |
| `%` | Bounce between open/close parentheses |
| `CTRL+d` | Page down (scroll cursor) |
| `CTRL+u` | Page down (scroll cursor) |
| `SHIFT+v` | Highlight full lines |

### Search

- Search all open buffers for text with `:bufdo`:

  ```vim
  :bufdo vimgrepadd threading % | copen
  ```

  - This will open the results in a quick fix window.
  - Use `CTRL+j`/`CTRL+k` to navigate between lines containing the search text.
  - Press enter to open the file to the matched line.

