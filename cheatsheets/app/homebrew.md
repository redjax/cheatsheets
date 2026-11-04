---
description: "[Homebrew](https://brew.sh) is a cross-platform package and dependency manager."
last_updated: "2026-11-04"
tags: ["app", "homebrew", "brew", "linux", "unix", "macos", "package-management"]
---

# Homebrew <!-- omit in toc -->

## About

[Homebrew](https://brew.sh/) is a package manager for macOS and Linux. Its `brew` CLI installs, updates, and removes packages, manages background services, and uses Brewfiles to record and restore package selections.

Terminology:

| Term     | Meaning                                                                              |
| -------- | ------------------------------------------------------------------------------------ |
| Formula  | Package definition, commonly used for command-line tools and libraries               |
| Cask     | Package definition for upstream-built binaries; commonly used for macOS applications |
| Tap      | Additional repository of formulae, casks, or Homebrew commands                       |
| Bottle   | Prebuilt package used instead of compiling a formula from source                     |
| Prefix   | Homebrew's installation directory                                                    |
| Keg      | Installation directory for one version of a formula                                  |
| Keg-only | Installed formula that is not symlinked into Homebrew's prefix                       |
| Brewfile | Ruby-based declaration of packages and other dependencies managed by `brew bundle`   |

📖 [`brew` Manual](https://docs.brew.sh/Manpage)

📖 [Homebrew Bundle and Brewfile](https://docs.brew.sh/Brew-Bundle-and-Brewfile)

## Usage

Replace `<formula>`, `<cask>`, `<package>`, and `<user/repo>` with actual names.

### Finding and inspecting packages

| Context     | Command                              | Action                                                                       |
| ----------- | ------------------------------------ | ---------------------------------------------------------------------------- |
| Search      | `brew search <text>`                 | Search formula and cask names                                                |
| Search      | `brew search '/<regex>/'`            | Search names using a regular expression                                      |
| Search      | `brew desc --search <text>`          | Search package names and descriptions                                        |
| Information | `brew info <package>`                | Show package details, installation status, dependencies, and caveats         |
| Information | `brew info --cask <cask>`            | Explicitly inspect a cask                                                    |
| Information | `brew home <package>`                | Open the package's homepage                                                  |
| Installed   | `brew list`                          | List installed formulae and casks                                            |
| Installed   | `brew list --formula`                | List installed formulae                                                      |
| Installed   | `brew list --cask`                   | List installed casks                                                         |
| Installed   | `brew list --versions`               | Show installed package versions                                              |
| Installed   | `brew list <formula>`                | List files installed by a formula                                            |
| Installed   | `brew leaves`                        | List formulae that are not dependencies of another installed formula or cask |
| Installed   | `brew leaves --installed-on-request` | List manually installed leaf formulae                                        |

### Installing and removing packages

| Context   | Command                                     | Action                                                                                        |
| --------- | ------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Install   | `brew install <formula>`                    | Install a formula and its dependencies                                                        |
| Install   | `brew install <formula1> <formula2>`        | Install multiple formulae                                                                     |
| Install   | `brew install --cask <cask>`                | Install a cask                                                                                |
| Reinstall | `brew reinstall <formula>`                  | Reinstall a formula                                                                           |
| Reinstall | `brew reinstall --cask <cask>`              | Reinstall a cask                                                                              |
| Remove    | `brew uninstall <formula>`                  | Uninstall a formula                                                                           |
| Remove    | `brew uninstall --cask <cask>`              | Uninstall a cask                                                                              |
| Remove    | `brew uninstall --cask --zap <cask>`        | Also remove associated files defined by the cask; may remove shared files or application data |
| Linking   | `brew unlink <formula>`                     | Remove a formula's symlinks without uninstalling it                                           |
| Linking   | `brew link <formula>`                       | Symlink an installed formula into Homebrew's prefix                                           |
| Linking   | `brew link --overwrite --dry-run <formula>` | Preview files that would be overwritten during linking                                        |

### Updating and maintenance

| Context      | Command                        | Action                                                                       |
| ------------ | ------------------------------ | ---------------------------------------------------------------------------- |
| Metadata     | `brew update`                  | Update Homebrew and package definitions; does not upgrade installed packages |
| Updates      | `brew outdated`                | List outdated installed packages                                             |
| Updates      | `brew upgrade`                 | Upgrade outdated formulae and eligible casks                                 |
| Updates      | `brew upgrade <formula>`       | Upgrade one formula                                                          |
| Updates      | `brew upgrade --cask <cask>`   | Upgrade one cask                                                             |
| Updates      | `brew upgrade --cask --greedy` | Include casks that auto-update or use an unversioned latest release          |
| Pinning      | `brew pin <formula>`           | Exclude an installed formula from normal upgrades                            |
| Pinning      | `brew list --pinned`           | List pinned formulae                                                         |
| Pinning      | `brew unpin <formula>`         | Allow a pinned formula to upgrade again                                      |
| Cleanup      | `brew cleanup --dry-run`       | Preview removal of old versions and stale cached files                       |
| Cleanup      | `brew cleanup`                 | Remove old versions and stale cached files                                   |
| Dependencies | `brew autoremove --dry-run`    | Preview removal of dependencies that are no longer needed                    |
| Dependencies | `brew autoremove`              | Uninstall dependencies that are no longer needed                             |

> [!NOTE]
> `brew update` updates Homebrew's knowledge of available software. `brew upgrade` updates installed software.
>
> Pinning does not create a portable version lock or make a pinned version available on another machine.

### Dependencies and taps

| Context      | Command                                       | Action                                        |
| ------------ | --------------------------------------------- | --------------------------------------------- |
| Dependencies | `brew deps <formula>`                         | List dependencies                             |
| Dependencies | `brew deps --tree <formula>`                  | Show dependencies as a tree                   |
| Dependencies | `brew deps --missing <formula>`               | Show missing dependencies                     |
| Dependents   | `brew uses --installed <formula>`             | List installed direct dependents              |
| Dependents   | `brew uses --installed --recursive <formula>` | List installed direct and indirect dependents |
| Taps         | `brew tap`                                    | List tapped repositories                      |
| Taps         | `brew tap <user/repo>`                        | Add a tap                                     |
| Taps         | `brew untap <user/repo>`                      | Remove a tap                                  |
| Taps         | `brew install <user/repo/formula>`            | Install a formula from a specific tap         |

### Services

For formulae that provide a service:

| Context     | Command                           | Action                                                  |
| ----------- | --------------------------------- | ------------------------------------------------------- |
| Status      | `brew services list`              | List managed services and their status                  |
| Information | `brew services info <formula>`    | Show service information                                |
| Start       | `brew services start <formula>`   | Start a service and register it to launch at login      |
| Run         | `brew services run <formula>`     | Run a service without registering it to launch at login |
| Stop        | `brew services stop <formula>`    | Stop a service and unregister it from login startup     |
| Restart     | `brew services restart <formula>` | Restart a service and register it to launch at login    |
| Cleanup     | `brew services cleanup`           | Remove unused service registrations                     |

These examples use user-level services without `sudo`. Services use `launchctl` on macOS and `systemctl` on Linux.

### Paths and diagnostics

| Context     | Command                   | Action                                                     |
| ----------- | ------------------------- | ---------------------------------------------------------- |
| Version     | `brew --version`          | Show Homebrew's version                                    |
| Help        | `brew help`               | Show general help                                          |
| Help        | `brew help <command>`     | Show help for a command                                    |
| Help        | `brew bundle --help`      | Show Bundle commands and options                           |
| Prefix      | `brew --prefix`           | Show Homebrew's installation prefix                        |
| Prefix      | `brew --prefix <formula>` | Show a formula's active installation path                  |
| Storage     | `brew --cellar`           | Show the Cellar directory                                  |
| Storage     | `brew --cache`            | Show the download cache directory                          |
| Shell       | `brew shellenv`           | Print shell commands that configure Homebrew's environment |
| Diagnostics | `brew doctor`             | Check for potential configuration problems                 |
| Diagnostics | `brew config`             | Show Homebrew and system configuration                     |
| Privacy     | `brew analytics`          | Show analytics status                                      |
| Privacy     | `brew analytics off`      | Disable Homebrew analytics                                 |

### Brewfiles and Bundle

`brew bundle` reads a Brewfile and installs its declared dependencies. `brew bundle install` is the explicit equivalent. By default, Bundle also upgrades outdated dependencies where supported.

Brew can find packages from other package managers installed on your system, including packages installed with `uv`, `mise`, and VS Code extensions. To export a clean list of only the apps you've installed with Homebrew, use:

```shell
brew bundle dump --formula --cask --tap --force
```

To export to a consistent location, add `--file=/path/to/Brewfile`:

```shell
brew bundle dump --formula --cask --tap --force --file="$HOME/dotfiles/Brewfile"
```

To show the output of `brew dump` without saving to a file, use an empty `--file=`:

```shell
brew bundle dump --formula --cask --tap --force --file=
```

#### Dumping installed packages

| Context         | Command                                                     | Action                                                    |
| --------------- | ----------------------------------------------------------- | --------------------------------------------------------- |
| Local dump      | `brew bundle dump`                                          | Create `./Brewfile`; refuse to overwrite an existing file |
| Local dump      | `brew bundle dump --force`                                  | Create or overwrite `./Brewfile`                          |
| Custom file     | `brew bundle dump --file="$HOME/dotfiles/Brewfile"`         | Dump to a specific file                                   |
| Custom file     | `brew bundle dump --file="$HOME/dotfiles/Brewfile" --force` | Overwrite a specific Brewfile                             |
| Global file     | `brew bundle dump --global`                                 | Dump to the configured global Brewfile                    |
| Global file     | `brew bundle dump --global --force`                         | Overwrite the configured global Brewfile                  |
| Standard output | `brew bundle dump --file=-`                                 | Print the generated Brewfile to stdout                    |
| Comments        | `brew bundle dump --force --no-describe`                    | Dump without package-description comments                 |
| Services        | `brew bundle dump --force --no-restart`                     | Omit service-restart declarations                         |
| Formulae        | `brew bundle dump --file=Brewfile.formulae --formula`       | Dump only formula entries                                 |
| Casks           | `brew bundle dump --file=Brewfile.casks --cask`             | Dump only cask entries                                    |
| Taps            | `brew bundle dump --file=Brewfile.taps --tap`               | Dump only tap entries                                     |
| Filtering       | `brew bundle dump --force --no-mas --no-vscode`             | Exclude Mac App Store apps and VS Code extensions         |

Type-specific dumps are useful for inspection or separate inventories. A formula-only dump may need tap declarations added before it can restore third-party packages elsewhere.

The global Brewfile location is selected from:

1. `$HOMEBREW_BUNDLE_FILE_GLOBAL`, if set.
2. `$XDG_CONFIG_HOME/homebrew/Brewfile`, if `$XDG_CONFIG_HOME` is set.
3. `~/.homebrew/Brewfile` or `~/.Brewfile`, according to Homebrew's global-file selection.

Use `--file` when you want an explicit, predictable location.

#### Restoring and checking packages

| Context        | Command                                        | Action                                                                  |
| -------------- | ---------------------------------------------- | ----------------------------------------------------------------------- |
| Local restore  | `brew bundle`                                  | Install and, by default, upgrade dependencies from `./Brewfile`         |
| Local restore  | `brew bundle install`                          | Explicit equivalent of `brew bundle`                                    |
| Custom restore | `brew bundle --file="$HOME/dotfiles/Brewfile"` | Restore from a specific file                                            |
| Global restore | `brew bundle --global`                         | Restore from the global Brewfile                                        |
| Avoid upgrades | `brew bundle --no-upgrade`                     | Install missing dependencies without explicitly upgrading outdated ones |
| Avoid upgrades | `brew bundle --global --no-upgrade`            | Apply the same behavior to the global Brewfile                          |
| Check          | `brew bundle check`                            | Check whether the Brewfile's dependencies are satisfied                 |
| Check          | `brew bundle check --verbose`                  | List unmet dependencies                                                 |
| Check          | `brew bundle check --no-upgrade`               | Check installation without treating outdated dependencies as unmet      |
| Inventory      | `brew bundle list --all`                       | List all dependency types declared by the Brewfile                      |
| Inventory      | `brew bundle list --formula`                   | List formula entries declared by the Brewfile                           |
| Inventory      | `brew bundle list --cask`                      | List cask entries declared by the Brewfile                              |
| Edit           | `brew bundle edit`                             | Open the Brewfile in an editor                                          |

`--no-upgrade` does not pin versions. Installing a missing package may still require upgrades to satisfy dependencies.

A restore normally leaves additional installed packages alone.

#### Removing packages not in a Brewfile

| Context   | Command                                                        | Action                                                                                          |
| --------- | -------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| Reconcile | `brew bundle cleanup --file="$HOME/dotfiles/Brewfile"`         | Review removal candidates through the interactive confirmation; decline to leave them installed |
| Reconcile | `brew bundle cleanup --file="$HOME/dotfiles/Brewfile" --force` | Remove supported dependencies not declared in the selected Brewfile without confirmation        |

Warning: Bundle cleanup is not the same as `brew cleanup`.

- `brew cleanup` removes stale files and old package versions.
- `brew bundle cleanup` can uninstall packages not declared in the Brewfile.
- A project-specific Brewfile is usually not a complete inventory of your machine.
- Current Bundle cleanup also reconciles Homebrew's explicit trust entries with those declared in the Brewfile.

## Examples

### Routine maintenance

```sh
# Refresh Homebrew and package definitions.
brew update

# Inspect available upgrades.
brew outdated

# Upgrade installed packages.
brew upgrade

# Preview optional cleanup.
brew autoremove --dry-run
brew cleanup --dry-run

# Remove unneeded dependencies and stale files after reviewing.
brew autoremove
brew cleanup
```

### Dump installed packages to a Brewfile

```sh
# Create a Brewfile in the current directory.
brew bundle dump

# Refresh it after installing or removing packages.
brew bundle dump --force

# Write a compact inventory without description comments.
brew bundle dump --force --no-describe

# Keep the snapshot in a specific directory.
mkdir -p "$HOME/dotfiles"
brew bundle dump --file="$HOME/dotfiles/Brewfile" --force

# Use Homebrew's global Brewfile location instead.
brew bundle dump --global --force

# Inspect the generated Brewfile without writing a file.
brew bundle dump --file=-
```

A dump records supported installed packages. Depending on the tools installed,
it may also include package types beyond formulae and casks, such as Mac App
Store apps or editor extensions.

Review the generated file before committing it.

### Restore packages from a Brewfile

Install Homebrew first, then:

```sh
# Restore from the current directory.
brew bundle install

# Restore from a specific file.
brew bundle install --file="$HOME/dotfiles/Brewfile"

# Install missing packages without explicitly upgrading existing ones.
brew bundle install --file="$HOME/dotfiles/Brewfile" --no-upgrade

# Restore from the global Brewfile.
brew bundle install --global

# Verify the declared dependencies.
brew bundle check --file="$HOME/dotfiles/Brewfile" --verbose
```

A Brewfile is not a full backup:

- It does not lock exact package versions.
- It does not back up application preferences, credentials, databases, or other user data.
- Restoring generally installs currently available versions.
- Packages still need to be available and compatible with the destination system.
- Normal restoration does not uninstall packages absent from the Brewfile.

### Keep a Brewfile in a dotfiles repository

```sh
# Record the current package selection.
brew bundle dump --file="$HOME/dotfiles/Brewfile" --force

# Review the changes.
git -C "$HOME/dotfiles" diff -- Brewfile

# Commit the snapshot.
git -C "$HOME/dotfiles" add Brewfile
git -C "$HOME/dotfiles" commit -m "Update Brewfile"

# On another machine, after cloning the repository and installing Homebrew:
brew bundle --file="$HOME/dotfiles/Brewfile"
```

Dumping with `--force` regenerates the file; it does not merge your hand-written
comments, conditional logic, or custom declarations. Keep a curated Brewfile
separate from generated snapshots if you want to preserve those edits.

### Write a Brewfile manually

```ruby
# Command-line tools.
brew "git"
brew "tig"
brew "jq"
brew "ripgrep"

# macOS applications.
if OS.mac?
  cask "firefox"
end
```

Install its dependencies:

```sh
brew bundle --file=./Brewfile
```

Brewfiles are evaluated as Ruby. Only run Brewfiles from sources you trust.

### Run commands with Brewfile dependencies

```sh
# Run a command in an environment configured for the Brewfile.
brew bundle exec -- git --version

# Open an interactive shell with the same environment.
brew bundle sh

# Check dependencies before running a command.
brew bundle exec --check -- git --version
```

This is especially useful for dependencies that are keg-only or otherwise not
on your normal shell's PATH.

## Troubleshooting

| Problem                                      | Check / command                                             | Notes                                                                                     |
| -------------------------------------------- | ----------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `brew: command not found`                    | Follow the installer's `brew shellenv` instructions         | Homebrew's executable directory may not be on PATH                                        |
| Installed command is not found               | `brew info <formula>`                                       | Read caveats; keg-only formulae may require PATH changes                                  |
| Wrong executable runs                        | `command -v <command>`                                      | Check whether another installation appears earlier on PATH                                |
| Need a package's install path                | `brew --prefix <formula>`                                   | Avoid hard-coding Homebrew's installation prefix                                          |
| General configuration warnings               | `brew doctor`                                               | Warnings are diagnostic hints; not every warning requires action if everything works      |
| Need details for a bug report                | `brew config`                                               | Includes Homebrew and system configuration                                                |
| Brewfile already exists                      | `brew bundle dump --force`                                  | Overwrites the file; review or commit existing edits first                                |
| Bundle cannot find the intended Brewfile     | `brew bundle --file="/full/path/to/Brewfile"`               | Use an explicit path rather than relying on the working directory or global configuration |
| Bundle check fails                           | `brew bundle check --verbose`                               | Shows unmet dependencies                                                                  |
| Bundle upgrades packages unexpectedly        | `brew bundle --no-upgrade`                                  | Disables explicit upgrades, but dependency installation may still require changes         |
| Restore does not match old versions          | Review Brewfile limitations                                 | Bundle has no exact-version lockfile                                                      |
| Generated Brewfile includes unwanted entries | Edit the file or use dump filters                           | For example, exclude Mac App Store apps with `--no-mas`                                   |
| Linking would overwrite files                | `brew link --overwrite --dry-run <formula>`                 | Inspect conflicts before using overwrite without the dry run                              |
| Unsure what Bundle cleanup removes           | Read `brew bundle --help` and inspect the selected Brewfile | Do not use a narrow project Brewfile as a machine-wide removal list                       |

## Links

- [Homebrew](https://brew.sh/)
- [Documentation](https://docs.brew.sh/)
- [`brew` Manual](https://docs.brew.sh/Manpage)
- [Homebrew Bundle and Brewfile](https://docs.brew.sh/Brew-Bundle-and-Brewfile)
- [Formula and Cask Search](https://formulae.brew.sh/)
- [Installation](https://docs.brew.sh/Installation)
- [Troubleshooting](https://docs.brew.sh/Troubleshooting)
- [FAQ](https://docs.brew.sh/FAQ)
- [Homebrew on GitHub](https://github.com/Homebrew/brew)
