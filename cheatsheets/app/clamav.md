---
description: "[ClamAV](https://www.clamav.net) is an open-source antivirus engine."
last_updated: "2026-05-24"
tags: ["app", "antivirus", "security"]
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
  - [Clamscan Cheatsheet](#clamscan-cheatsheet)
  - [Freshclam Cheatsheet](#freshclam-cheatsheet)
  - [Clamdscan Cheatsheet](#clamdscan-cheatsheet)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Clamav <!-- omit in toc -->

## About

[ClamAV](https://www.clamav.net) is an open-source antivirus engine.

## Usage

### Clamscan Cheatsheet

| command/arg                  | description                                    |
| ---------------------------- | ---------------------------------------------- |
| `clamscan FILE`              | Scan a single file                             |
| `clamscan DIR`               | Scan a directory (non-recursive by default)    |
| `-r`, `--recursive`          | Recursively scan directories                   |
| `-i`, `--infected`           | Only show infected files                       |
| `--remove=yes`               | Remove infected files automatically            |
| `--move=/path/quarantine`    | Move infected files to a quarantine directory  |
| `--copy=/path/quarantine`    | Copy infected files to a quarantine directory  |
| `-l FILE`, `--log=FILE`      | Save scan output to a log file                 |
| `--quiet`                    | Only display errors                            |
| `--verbose`                  | Show detailed scan output                      |
| `--stdout`                   | Write all output to stdout                     |
| `--no-summary`               | Do not display scan summary                    |
| `--bell`                     | Ring terminal bell when malware is found       |
| `--max-filesize=SIZE`        | Skip files larger than specified size          |
| `--max-scansize=SIZE`        | Limit total data scanned inside archives       |
| `--max-files=NUMBER`         | Limit number of files scanned inside archives  |
| `--database=/path/db`        | Use a custom virus database path               |
| `--official-db-only=yes`     | Only load official ClamAV signatures           |
| `--tempdir=/path/tmp`        | Use a custom temp directory                    |
| `--leave-temps`              | Keep temporary extracted files                 |
| `--scan-archive=yes`         | Scan archive files (zip, tar, rar, etc.)       |
| `--scan-mail=yes`            | Scan email files                               |
| `--scan-pdf=yes`             | Scan PDF files                                 |
| `--scan-html=yes`            | Scan HTML files                                |
| `--phishing-scans=yes`       | Enable phishing detection                      |
| `--heuristic-alerts=yes`     | Enable heuristic malware detection             |
| `--detect-pua=yes`           | Detect potentially unwanted applications (PUA) |
| `--exclude=REGEX`            | Exclude files matching regex                   |
| `--exclude-dir=REGEX`        | Exclude directories matching regex             |
| `--include=REGEX`            | Only scan files matching regex                 |
| `--follow-dir-symlinks=yes`  | Follow symbolic links to directories           |
| `--follow-file-symlinks=yes` | Follow symbolic links to files                 |
| `--cross-fs=yes`             | Scan across mounted filesystems                |
| `--bytecode-timeout=N`       | Set bytecode execution timeout                 |
| `--alert-encrypted=yes`      | Alert on encrypted archives/documents          |
| `--nocerts`                  | Disable certificate verification in PE files   |
| `--disable-cache`            | Disable scan cache                             |
| `--file-list=FILE`           | Read list of files to scan from a file         |
| `-h`, `--help`               | Display help                                   |
| `-V`, `--version`            | Display ClamAV version                         |

### Freshclam Cheatsheet

`freshclam` is a utility to update definitions for the AV engine.

| command/arg                    | description                       |
| ------------------------------ | --------------------------------- |
| `freshclam`                    | Update virus definitions          |
| `freshclam -v`                 | Verbose update output             |
| `freshclam --stdout`           | Send output to stdout             |
| `freshclam --quiet`            | Suppress normal output            |
| `freshclam --check=N`          | Check for updates N times per day |
| `freshclam --daemon`           | Run as a background daemon        |
| `freshclam --log=FILE`         | Write update logs to file         |
| `freshclam --datadir=/path/db` | Use custom database directory     |

### Clamdscan Cheatsheet

`clamdscan` uses the daemon for scanning, and is the preferred command over `clamscan`. Scans are faster, can be continuous or scheduled, and the commands are generally simpler.

| command/arg          | description                              |
| -------------------- | ---------------------------------------- |
| `clamdscan FILE`     | Scan file using running `clamd` daemon   |
| `clamdscan DIR`      | Scan directory using daemon              |
| `--multiscan`        | Use multiple threads for faster scanning |
| `--fdpass`           | Pass file descriptors to daemon          |
| `--stream`           | Stream file contents to daemon           |
| `--config-file=FILE` | Use custom `clamd.conf`                  |
| `--move=DIR`         | Move infected files                      |
| `--remove=yes`       | Remove infected files                    |
| `--infected`         | Only show infected files                 |
| `--quiet`            | Only show errors                         |

## Examples

Update virus definitions

```shell
sudo freshclam
```

Scan current directory

```shell
clamscan .
```

Recursively scan a directory

```shell
clamscan -r /home/user/downloads
```

Only show infected files

```shell
clamscan -r -i /home/user
```

Scan entire Linux system

```shell
sudo clamscan -r /
```

Save scan results to a log file

```shell
clamscan -r /home/user --log=/var/log/clamav/scan.log
```

Move infected files to quarantine

```shell
clamscan -r --move=/quarantine /home/user
```

Remove infected files automatically

```shell
clamscan -r --remove=yes /home/user
```

Scan files listed in a text file

```shell
clamscan --file-list=scan_targets.txt
```

Scan and only display infected files

```shell
clamscan -r -i --bell /home/user
```

Exclude specific directories

```shell
clamscan -r / --exclude-dir="^/sys|^/proc|^/dev"
```

Detect potentially unwanted applications (PUA)

```shell
clamscan -r --detect-pua=yes /home/user
```

Use a Custom virus database directory

```shell
clamscan --database=/opt/clamav/db suspicious_file.exe
```

Faster scanning with clamdscan

```shell
clamdscan --multiscan /home/user
```

Scan a file from standard input

```shell
cat suspicious_file | clamscan -
```

Verbose scan output

```shell
clamscan -r --verbose /home/user
```

Quiet scan (errors only)

```shell
clamscan -r --quiet /home/user
```

## Troubleshooting

## Links

- [ClamAV home](https://clamav.net)
- [ClamAV docs](https://docs.clamav.net/)
- [ClamAV scanning docs](https://docs.clamav.net/manual/Usage/Scanning.html)
