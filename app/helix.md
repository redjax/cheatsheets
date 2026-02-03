---
description: "Terminal text editor like Neovim, but with a bunch of built-in functionality instead of relying on plugins."
last_updated: "2026-02-02"
tags: ["helix", "app", "cli", "tui"]
---

# Helix Editor

[https://helix-editor.com](https://helix-editor.com)

## Usage

Run `hx --tutor` to see Helix's built-in tutorial.

- `hjkl`/arrow keys: Move cursor around.
- `d`: Delete the character under the cursor.
- `o`: Start a new line below the current line and enter insert mode.
- `O`: Start a new line above the current line and enter insert mode.
- `a`: Append to the current selection.
- `I`: Enter insert mode at the first non-whitespace character on the current line.
- `A`: Enter insert mode at the end of the current line.
- `w`: Jump cursor to beginning of next word.
  - Start a selection with `w`, then use `e`/`b` to highlight text.
  - Use capitals `W`/`E`/`B` to select whole words, including ones separated by characters like `-`.
- `e`: Jump cursor to end of the current word.
- `b`: Jump backward to beginning of current word.
- `c`: Change the current selection.
  - Change a single character, or use `w` (or `e` or `b`) to select it and then use `c` to change/replace the whole word.
  - You can highlight as much as you want and use `c` to replace the text.
- `C`: Add a new cursor on the next suitable line, skipping empty lines.
  - `ALT + SHIFT + C`: Add a new cursor on the previous suitable line, skipping newlines.
  - `,`: Return to 1 cursor.
- `#w`/`#e`/`#b`: Move cursor forward/backwards n number of words, where n is the value you used for `#`.
  - Examples:
    - `2w`: Move cursor forward 2 words.
    - `3e`: Move cursor to end of the 3rd word and forward.
    - `2b`: Move 2 words backwards.
- `v`: Enter select mode (`v` or `ESC` to exit).
  - Combine with a number and direction to select by word, i.e. `v2w`.
- `%`: Select all text in the file.
- `x`: Select whole line. Press again to select multiple lines.
- `#x`: Select n lines, where n is the value you used for `#`.
- `X`: Select without extending to next line.
- `;`: Collapse selections (un-highlight).
- `u`: Undo last action. Can be repeated.
- `U`: Redo last action. Can be repeated.
- `y`: Yank selection.
- `Space + y`: Yank to system clipboard.
- `p`: Paste yanked selection.
- `Space + p`: Paste from the system keyboard.
- `/`: Start searching current buffer.
  - `n`: Go to next search match.
  - `P`: Go to previous search match.
- `s`: Select all similar words in current selection.
  - For example, highlight a line with `x`, then press `s` to search for a word that appears multiple times in that line.
  - You can use regex in a search, i.e. by selecting a line with `x`, pressing `s` to start a search.
    - For example, to search for all instances of 2 or more `x` characters (`x` could be a space character): `xx+`.
  - Here is how to do a search and replace on the whole file:
    - Select all of the text in the file with `%`.
    - Press `s` to start a search, type the word/term/phrase you want to replace.
    - Use `i`/`a` to enter insert mode and type, or `c` to change all instances.
- `&`: Align contents.
  - Example: Align numbers
    - Say you have 4 numbered lines, like:

        ```text
        1)  lorem
        2)  ipsum
        3)  dolor
        4)   sit
        ```

    - Place the cursor in the whitespace after the "97) ".
    - Press `C` 4 times to add a new cursor on the lines below.
    - Press `W` to select the numbers and brackets on each line.
    - Press `&` to align all of the lines.
  - Example: Align Markdown table
    - Say you have a Markdown table like:

      ```markdown
          | FRUIT   | AMOUNT |
          |---------|--------|
      | Apples  | 8      |
          | Bananas | 6      |
      | Oranges | 3      |
          | Donuts  | 4      |
      ```

    - Select the whole table, i.e. by putting the cursor at the start of the table and typing `6x`.
    - Press `ALT-s` to create a cursor on each line.
    - Press `&` to align the table.
- `f<char>`: Select up to and including a character, i.e. `f[` or `f{`.
- `F<char>`: Select backwards up to and including a character, i.e. `F[` or `F]`.
- `t<char>`: Select up to a character, but don't select the character itself. i.e. `t[` `t]`.
- `T<char>`: Select backwards up to a character, but don't select the character itself. i.e. `T[` `T]`.
- `r<char>`: Replace all selected characters with given character.
  - Example: You have a Markdown table using `=` instead of `-` for dividers:

    ```markdown
    | Month | Days |
    |=======|------|
    | Jan   | 31   |
    | Feb   | 28   |
    | Mar   | 31   |
    | ...   | ...  |
    ```

    - Put your cursor on the first `=`.
    - Press `t|` to select that character all the way to the `|`, without selecting `|`.
    - Press `r-` to replace all of the selected characters with `-`.
- `R`: Replace selection with yanked text.
  - Select text and yank with `y`.
  - Select the text you want to replace and press `R`.
  - Example: Say you have text split over multiple lines like this:

    ```text
    This sentence
    is spilling over
    onto other
    lines.
    ```

    - Select the 4 lines by putting the cursor on the first line (`This sentence`) and pressing x 4 times.
    - Press `J` to join all of the lines onto 1.
- `>`: Indent line.
- `<`: Unindent line.
- `CTRL + a`: Increment a selected number.
  - i.e. `1)` -> `CTRL+a` -> `2)`.
- `CTRL + x`: Decrement a selected number.
  - i.e. `2)` -> `CTRL+x` -> `1)`.
- `"<char>`: Yank text to a named register, denoted with `<char>`. Paste later with `<char> R`.
  - Example: Say you have this text:

    ```text
    I like watermelons and bananas because my favorite fruits are watermelons and pineapples.
    ```

    - Select the text `watermelons` and yank.
    - Select the text `bananas` and yank it to a named register, i.e. `"b` then `y`.
      - This stores the value `bananas` in a register named `b`.
    - Select a word, i.e. `mangoes`, and press `R` to paste the yanked yext (`watermelons`).
    - Select another word, i.e. `pineapples`, select the `b` register and press `R` to replace it with `bananas`.
      - i.e. `b R` after selecting the word `pineapples`.
- `*`: Jump to next entry in selection search.
  - Highlight a word/phrase.
  - Press `*` to copy it to a register.
  - Press `n` and `N` to jump between instances of the selected text.
- `n`: While in select mode (`v`), select the next matching word/text.
  - Highlight desired text that repeats on other lines.
  - Highlight the word and press `*` to store it.
  - Type `v` to enter select mode and hit `n` to select other instances.
  - Use `c` or `r` to change the text.
  - I.e. change all instances of `bat` to `cat`.
- `N`: While in select mode (`v`), select the previous matching word/text (functions the same as `n`).
- `CTRL + s`: Save a jump location.
  - You can use `CTRL + i` and `CTRL + o` to cycle forwards/backwards through jump locations.
    - `CTRL + i`: Move forward to next jump location.
    - `CTRL + o`: Move backwards to previous jump location.
- `(`/`)`: Cycle selection forwards/backwards.
  - Select multiple lines, i.e. with `x` `x` (or `2x`).
  - Press `s` and type a word to highlight in the selection.
  - Use `(` to jump backwards to the previous selection.
  - Use `)` to jump forwards to the next selection.
- `~`: Switch case of selected letters.
- `` ` ``: Set all selected characters to lowercase.
- `` Alt + ` ``: Set all selected characters to uppercase.
- `CTRL + C`: Comment/uncomment a line.
