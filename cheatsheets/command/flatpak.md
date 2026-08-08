---
description: "[Flatpak](https://flatpak.org) runs applications in sandboxed environments."
last_updated: "{{last_update}}"
tags: ["command", "flatpak", "linux"]
last_updated: "2026-08-08"
---

# Flatpak <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [About](#about)
  - [Installing Flatpak](#installing-flatpak)
  - [Enable Flathub](#enable-flathub)
  - [Quick Reference](#quick-reference)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

## About

[Flatpak](https://flatpak.org) runs applications in sandboxed environments. You can find new flatpaks on [Flathub](https://flathub.org). Most Linux app stores like Gnome's Software and KDE's Discover will offer Flatpak installation once it's installed.

### Installing Flatpak

| Distro / Family         | Command                                                  | Notes                                              |
| ----------------------- | -------------------------------------------------------- | -------------------------------------------------- |
| **Debian / Ubuntu**     | `sudo apt install flatpak`                               | Install Flatpak from the distro repositories       |
| **Fedora / RHEL-based** | `sudo dnf install flatpak`                               | Works on Fedora and other modern DNF-based systems |
| **Arch Linux**          | `sudo pacman -S flatpak`                                 | Install from the official Arch repositories        |
| **openSUSE**            | `sudo zypper install flatpak`                            | Install using Zypper                               |
| **Alpine Linux**        | `sudo apk add flatpak`                                   | Install from Alpine's package repositories         |
| **Gentoo**              | `sudo emerge --ask sys-apps/flatpak`                     | Install through Portage                            |
| **NixOS**               | `sudo nix-channel --update && sudo nixos-rebuild switch` | Enable Flatpak in your NixOS configuration first   |
| **Void Linux**          | `sudo xbps-install -S flatpak`                           | Install from the Void package repositories         |
| **Solus**               | `sudo eopkg install flatpak`                             | Install using Solus' package manager               |
| **Mageia**              | `sudo urpmi flatpak`                                     | Install from Mageia repositories                   |

### Enable Flathub

After installing Flatpak, Flathub is commonly added as the application source:

```bash
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
```

Verify that it was added:

```bash
flatpak remotes
```

Then you can install applications, for example:

```bash
flatpak install flathub org.mozilla.firefox
```

### Quick Reference

| Family          | Package manager | Install Flatpak               |
| --------------- | --------------- | ----------------------------- |
| Debian / Ubuntu | `apt`           | `sudo apt install flatpak`    |
| Fedora / RHEL   | `dnf`           | `sudo dnf install flatpak`    |
| Arch            | `pacman`        | `sudo pacman -S flatpak`      |
| openSUSE        | `zypper`        | `sudo zypper install flatpak` |
| Alpine          | `apk`           | `sudo apk add flatpak`        |

## Usage

| Command                           | Example                                                                                      | Description                                             |
| --------------------------------- | -------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| `flatpak install`                 | `flatpak install flathub org.mozilla.firefox`                                                | Install an application from a remote repository         |
| `flatpak uninstall`               | `flatpak uninstall org.mozilla.firefox`                                                      | Uninstall an application                                |
| `flatpak update`                  | `flatpak update`                                                                             | Update installed applications and runtimes              |
| `flatpak run`                     | `flatpak run org.mozilla.firefox`                                                            | Launch an installed application                         |
| `flatpak list`                    | `flatpak list`                                                                               | List installed applications and runtimes                |
| `flatpak search`                  | `flatpak search firefox`                                                                     | Search available applications                           |
| `flatpak info`                    | `flatpak info org.mozilla.firefox`                                                           | Show information about an installed application         |
| `flatpak remote-list`             | `flatpak remote-list`                                                                        | List configured repositories/remotes                    |
| `flatpak remote-add`              | `flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo` | Add a Flatpak repository                                |
| `flatpak remote-delete`           | `flatpak remote-delete flathub`                                                              | Remove a configured repository                          |
| `flatpak repair`                  | `flatpak repair`                                                                             | Check and repair the local Flatpak installation         |
| `flatpak uninstall --unused`      | `flatpak uninstall --unused`                                                                 | Remove unused runtimes and other unused objects         |
| `flatpak history`                 | `flatpak history`                                                                            | Show Flatpak installation/update history                |
| `flatpak override`                | `flatpak override --user --filesystem=home org.example.App`                                  | Modify an application's sandbox permissions             |
| `flatpak override --show`         | `flatpak override --show org.example.App`                                                    | Show permission overrides for an application            |
| `flatpak permission-show`         | `flatpak permission-show org.example.App`                                                    | Show application permissions                            |
| `flatpak kill`                    | `flatpak kill org.example.App`                                                               | Stop a running Flatpak application                      |
| `flatpak ps`                      | `flatpak ps`                                                                                 | List running Flatpak applications                       |
| `flatpak make-current`            | `flatpak make-current org.example.App stable`                                                | Select the current branch/version of an application     |
| `flatpak mask`                    | `flatpak mask org.example.App`                                                               | Prevent an application or runtime from being updated    |
| `flatpak pin`                     | `flatpak pin org.gnome.Platform//48`                                                         | Keep a runtime installed even if it is otherwise unused |
| `flatpak uninstall --delete-data` | `flatpak uninstall --delete-data org.example.App`                                            | Uninstall an app and delete its stored user data        |
| `flatpak config`                  | `flatpak config --list`                                                                      | View or modify Flatpak configuration                    |
| `flatpak remote-info`             | `flatpak remote-info flathub org.mozilla.firefox`                                            | Show information about an app available from a remote   |
| `flatpak install --user`          | `flatpak install --user flathub org.example.App`                                             | Install an application for the current user only        |
| `flatpak uninstall --user`        | `flatpak uninstall --user org.example.App`                                                   | Uninstall a per-user application                        |
| `flatpak list --app`              | `flatpak list --app`                                                                         | List only installed applications                        |
| `flatpak list --runtime`          | `flatpak list --runtime`                                                                     | List only installed runtimes                            |

## Examples

| Cmd                     | Example                                              | Description                                                  |
| ----------------------- | ---------------------------------------------------- | ------------------------------------------------------------ |
| Install from Flathub    | `flatpak install flathub APP_ID`                     | Install an app using its Flatpak application ID              |
| Run by ID               | `flatpak run APP_ID`                                 | Launch an application without needing its desktop menu entry |
| Update everything       | `flatpak update`                                     | Update all installed Flatpaks                                |
| Remove leftovers        | `flatpak uninstall --unused`                         | Clean up runtimes no longer needed                           |
| Find an app ID          | `flatpak search SEARCH_TERM`                         | Search for the exact application ID                          |
| Inspect permissions     | `flatpak info APP_ID`                                | View application details, including runtime information      |
| Grant filesystem access | `flatpak override --user --filesystem=PATH APP_ID`   | Give an app access to a filesystem path                      |
| Revoke an override      | `flatpak override --user --nofilesystem=PATH APP_ID` | Remove a previously granted filesystem override              |
| Reset overrides         | `flatpak override --user --reset APP_ID`             | Reset custom permissions for an application                  |

## Troubleshooting

## Links

- [Flatseal: Flatpak permission manager](https://github.com/tchx84/Flatseal)
