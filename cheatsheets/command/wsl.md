---
description: "WSL is the Windows CLI package manager, like apt on Ubuntu or DNF on Fedora."
last_updated: "{{last_update}}"
tags: ["command", "windows", "wsl"]
last_updated: "2026-03-16"
---

# Windows Subsystem for Linux (WSL) <!-- omit in toc -->

## About

WSL is the Windows CLI package manager, like apt on Ubuntu or DNF on Fedora.

## Usage

### Commands Cheatsheet

| Command | Description |
| ------- | ----------- |
| `wsl --list --online` | List available WSL distributions. |
| `wsl --install -d <distro-name>` | Install a WSL distribution, i.e. `Debian` or `Ubuntu`. |
| `wsl -t <distro-name>` | Terminate/stop a running WSL distribution. |
| `wsl --unregister <distro-name>` | Remove/uninstall a WSL distribution. |
| `wsl --export <distro-name> c:\Path\To\distroname-backup.tar` | Create a backup of a WSL distribution. |
| `wsl --import <distro-name> c:\Path\To\Install\<distro-name> c:\Path\To\distroname-backup.tar` | Restore a WSL backup. |
| `wsl --shutdown` | Shutdown WSL service host, killing all running VMs. |


### List WSL distributions available for install

```shell
wsl --list --online
```

### Add a WSL distribution

```shell
wsl --install -d <distro-name>
```

For example, to install Debian:

```shell
wsl --install -d Debian
```

### Remove a WSL distribution

```shell
wsl --unregister <distro-name>
```

### Shutdown/Stop a distribution

To stop a running distribution:

```shell
wsl -t <distro-name>
```

#### Shutdown all distributions

To fully shut down all WSL distributions:

```shell
wsl --shutdown
```

### Backup & Restore

#### Backup WSL distro

You can use the `--export` flag to create a `.tar` backup of a WSL distribution. For example, to backup a distribution named `debian` to the user's `Documents\wsl_backup` directory:

```shell
wsl --export debian c:\Users\username\documents\wsl_backup\debian.tar
```

#### Restore WSL distro from backup

Restoring a distro requires 3 parameters:

- distribution: The type of Linux (Debian, Ubuntu, etc)
- install location: The place to restore the WSL machine to, i.e. `c:\wsl`
- backup path: Path to the WSL backup to restore

For example, to restore a distro named `debian` from a user's `Documents\wsl_backup` directory:

```shell
wsl --import debian c:\wsl\Debian c:\users\username\documents\wsl_backup\debian.tar
```

### Enable systemd in WSL

Edit `/etc/wsl.conf` inside the WSL container and add this:

```conf
[boot]
systemd=true
```

### Enable FUSE fs in WSL

On the Windows host, edit `%USERPROFILE%\.wslconfig` and add:

```conf
[wsl2]
kernelCommandLine = "fuse.enable=1"
```

### Use host's git credential manager

In your WSL distribution, run:

```shell
git config --global credential.helper "/mnt/c/Program\ Files/Git/mingw64/libexec/git-core/git-credential-wincred.exe"
```

If Git was installed to the user's AppData directory, use this command (replacing `<username>` with the Windows user where Git was installed):

```shell
git config --global credential.helper "/mnt/c/Users/<username>/AppData/Local/Programs/Git/mingw64/bin/git-credential-manager.exe"
```

If your git remote is Azure DevOps, also run:

```shell
git config --global credential.https://dev/azure.com.useHttpPath true
```

### Use git in WSL without Windows GCM

To use Git in WSL without calling the Windows host's Git Credential Manager:

```shell
git config --global credential.helper ""
git config --global credential.helper cache
git config --global credential.cacheTimeout 86400
```

If you want to enable saving credentials to a file, run:

```shell
git config --global credential.helper store
```

## Examples

## Troubleshooting

### Fix signature mismatches in Azure libraries

- [Enable systemd in WSL](#enable-systemd-in-wsl)
- Add the following line to the `[boot]` section:

```conf
[boot]
command="ntpdate ntp.ubuntu.com"
```

### Fix ping socket operation not permitted

Example error:

```shell
ping: socktype: SOCK_RAW
ping: socket: Operation not permitted
ping: => missing cap_net raw+p capability or setuid?
```

To fix, edit `%USERPROFILE%\.wslconfig` on the Windows host and add this:

```conf
[wsl2]
kernelCommandLine = sysctl.net.ipv4.ping_group_range=\"0 2147483647\"
```

### Kill frozen or suspended WSL

Sometimes WSL will completely freeze up and stop working, especially if a WSL distribution was running when a machine went to sleep. To fix, you can kill and restart the WSL service with:

```shell
taskkill /f /im wslservice.exe
```

Then relaunch with `wsl` or `wsl -d <distro-name>`

## Links
