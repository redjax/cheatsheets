---
description: "Terminal-based IDE, next-gen Vi(m)."
last_updated: "2026-02-20"
tags: ["neovim", "app", "cli", "tui"]
---

# Neovim <!-- omit in toc -->

[https://neovim.io/](https://neovim.io/)

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
  - [CLI args](#cli-args)

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
