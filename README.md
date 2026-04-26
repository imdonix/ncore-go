# nCore REST API (Unofficial)

A Go library and REST API for interacting with torrents from [ncore.pro](https://ncore.pro).

> **Disclaimer**: This is an unofficial project. It is not affiliated with ncore.pro.

## Features

- Search torrents with various filters (type, category, sort order)
- Get detailed torrent information by ID
- Retrieve Hit & Run activity data
- Get recommended torrents
- Download torrent files directly
- Cookie-based authentication session management
- Optional 2FA support

## Library Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/imdonix/ncore-go/pkg/ncore"
)

func main() {
    client, _ := ncore.NewClient(15*time.Second, nil)
    
    cookies, err := client.Login("username", "password", "")
    if err != nil {
        panic(err)
    }
    fmt.Println("Logged in, cookies:", cookies)

    result, _ := client.Search("movie name", ncore.TypeAllOwn, ncore.WhereName, ncore.SortSeeders, ncore.SeqDesc, 1)
    for _, t := range result.Torrents {
        fmt.Printf("%s - %s (%d seeders)\n", t.Name, t.Size, t.Seeders)
    }
}
```

## REST API

The server runs on port 8080 by default.

### Authentication

**Login**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "your_user", "password": "your_pass", "2factor": "123456"}'
```

Returns auth cookies for use in subsequent requests.

**Logout**
```bash
curl -X POST http://localhost:8080/logout \
  -H "X-Ncore-nick: your_nick" \
  -H "X-Ncore-pass: your_pass" \
  -H "X-Ncore-stilus: default" \
  -H "X-Ncore-nyelv: hu" \
  -H "X-Ncore-PHPSESSID: your_session"
```

### Endpoints

All endpoints (except `/login`) require authentication cookies passed as headers.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/login` | Authenticate and get cookies |
| `POST` | `/search` | Search torrents |
| `GET` | `/torrent/:id` | Get torrent details |
| `GET` | `/torrent/:id/download` | Download torrent file |
| `GET` | `/activity` | Get Hit & Run list |
| `GET` | `/recommended` | Get recommended torrents |
| `POST` | `/logout` | Clear session |

### Search Request

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -H "X-Ncore-nick: your_nick" \
  -H "X-Ncore-pass: your_pass" \
  -H "X-Ncore-stilus: default" \
  -H "X-Ncore-nyelv: hu" \
  -H "X-Ncore-PHPSESSID: your_session" \
  -d '{
    "pattern": "lord of the rings",
    "type": "all_own",
    "where": "name",
    "sort_by": "seeders",
    "sort_order": "DESC",
    "page": 1
  }'
```

**Search Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `pattern` | string | Search query |
| `type` | string | Torrent category |
| `where` | string | Search field (name, description, imdb, label) |
| `sort_by` | string | Sort field (name, fid, size, times_completed, seeders, leechers) |
| `sort_order` | string | Sort order (ASC, DESC) |
| `page` | int | Page number |

**Categories**: `xvid_hun`, `xvid`, `dvd_hun`, `dvd`, `dvd9_hun`, `dvd9`, `hd_hun`, `hd`, `xvidser_hun`, `xvidser`, `dvdser_hun`, `dvdser`, `hdser_hun`, `hdser`, `mp3_hun`, `mp3`, `lossless_hun`, `lossless`, `clip`, `game_iso`, `game_rip`, `console`, `ebook_hun`, `ebook`, `iso`, `misc`, `mobil`, `xxx_imageset`, `xxx_xvid`, `xxx_dvd`, `xxx_hd`, `all_own`

### Download Torrent

```bash
curl -X GET "http://localhost:8080/torrent/123456/download" \
  -H "X-Ncore-nick: your_nick" \
  -H "X-Ncore-pass: your_pass" \
  -H "X-Ncore-stilus: default" \
  -H "X-Ncore-nyelv: hu" \
  -H "X-Ncore-PHPSESSID: your_session" \
  -o torrent.torrent
```

## Docker

```bash
docker build -t ncore-go .
docker run -p 8080:8080 ncore-go
```

## Build

```bash
make build    # Build binary to bin/ncore
make test    # Run tests
make format  # Format code
make docker  # Build Docker image
```
