# Cheatsheets

My personal `~/.cheatsheets`. These are `man`-style documents written in Markdown, which I can easily open in a terminal editor or pager to view/search.

Cheatsheets are in Markdown format, and can be used cross-platform, version controlled, and opened in any number of editors or viewers.

## Index

Check the [auto-generated index file](./INDEX.md) for a map of the cheatsheets.

## Setup

Clone the repository to a directory, i.e. `~/.cheatsheets`:

```shell
git clone git@github.com:redjax/cheatsheets.git ~/git/cheatsheets
```

Create a symlink of the [`cheatsheets/` directory](./cheatsheets/):

```shell
ln -s $HOME/git/cheatsheets $HOME/.cheatsheets
```

## CLI Tool

The `app/` directory contains a Go CLI tool for managing cheatsheets.

The CLI tool clones a "working copy" of the repository to a location defined in your config (default: `~/.local/share/cheatsheets` on Linux/Mac, `%APPDATA%\cheatsheets` on Windows). It checks out a branch named `working` so any changes you made can be safely merged into main when ready.

You can `cd` to the `chtsht`-managed repository using `chtsht cd`, and synchronize changes with `chtsht sync`.

## Install CLI

#### Quick Install (recommended)

- Linux / macOS (bash):

```shell
curl -fsSL https://raw.githubusercontent.com/redjax/cheatsheets/main/.scripts/install.sh | bash
```

- Windows (PowerShell):

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/redjax/cheatsheets/main/.scripts/install.ps1)))
```

> [!TIP]
> Add `-s -- --auto` (bash) or `-Auto` (PowerShell) to skip confirmation prompts.

The install script downloads the latest release from GitHub and installs to:

- Linux / macOS: `~/.local/bin/chtsht`
- Windows: `%LOCALAPPDATA%\chtsht\chtsht.exe` (added to user PATH automatically)

#### Manual Install

Check the [releases page](https://github.com/redjax/cheatsheets/releases), look for releases with `v0.0.0` tags (the Go app). Cheatsheet releases are tagged like `Cheatsheets <date>-<shorthash>`.

Download the zip for your platform, extract the binary, and place it somewhere on your `$PATH`.

### Build from source

```shell
cd app/
go build -o chtsht ./cmd/chtsht/
# Move the chtsht binary to a location on your $PATH or use ./chtsht
```

### Configuration

Copy `config.yml` to `config.local.yml` and edit with your settings:

```yaml
git:
  repo_url: "https://github.com/your-username/cheatsheets.git"
  token: "your-github-token"  # For push access
  auto_branch: true
  working_branch: "working"
```

### Usage

> [!NOTE]
> Use `chtsht -h` to see help menu

View available cheatsheets:

```shell
chtsht list
chtsht list --type app
```

Edit a cheatsheet:

```shell
chtsht edit app/neovim
```

Create new cheatsheet:

```shell
chtsht new --type language --name rust
```

Delete cheatsheet:

```shell
chtsht delete app/helix
```

### Git Workflow

The tool uses a working branch workflow. Edit on the `working` branch, merge to `main` when ready.

Check status:

```shell
chtsht repo status
```

Stage, commit, push changes:

```shell
chtsht repo stage --all
chtsht repo commit -m "message"
chtsht repo push
```

Or sync in one command:

```shell
chtsht repo sync -m "message"
```

Pull latest changes:

```shell
chtsht repo pull
```

Merge working branch to main:

```shell
chtsht repo merge-to-main
```

Branch management:

```shell
chtsht repo branch list
chtsht repo branch ensure  # Create/switch to working branch
chtsht repo branch switch main
```
