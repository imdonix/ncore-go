# Ncore (Go)

## Introduction

This repository provides a Go API to manage torrents from ncore.pro.

## Install

```bash
go get github.com/imdonix/ncore-go/pkg/ncore
```

## Structure

- `pkg/ncore`: The public library logic.
- `cmd/example`: An example application demonstrating how to use the library.

## Features
- Search torrents with various filters.
- Get torrent details by ID.
- Get torrents from activity (Hit & Run).
- Get recommended torrents.
- Download torrent files.
- Cookie-based authentication.
- Supports 2FA.
