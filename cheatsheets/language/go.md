---
description: "[Go](https://go.dev) is a programming language with a friendly developer experience & efficient binaries."
last_updated: "{{last_update}}"
tags: ["language", ]
last_updated: "2026-06-02"
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
  - [Download & install Go](#download-install-go)
  - [Update all packages in current module](#update-all-packages-in-current-module)
- [Usage](#usage)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Go <!-- omit in toc -->

## About

[Go](https://go.dev) is a programming language with a friendly developer experience & efficient binaries.

### Download & install Go

Find the latest version of Go at the [Go releases dashboard](https://go.dev/dl), or using the [https://go.dev/VERSION?m=text URL](https://go.dev/VERSION?m=text).

The `/VERSION` endpoint can be scripted, too, i.e. to install Go on a 64-bit Linux machine:

```shell
wget "https://go.dev/dl/$(curl 'https://go.dev/VERSION?m=text').linux-amd64.tar.gz" && sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz
```

### Update all packages in current module

```shell
go get -u ./...
go mod tidy
```

## Usage

## Examples

## Troubleshooting

## Links
