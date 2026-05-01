# nCore REST API (Unofficial)

A Go library and REST API for interacting with torrents from [ncore.pro](https://ncore.pro).

> **Disclaimer**: This is an unofficial project. It is not affiliated with ncore.pro.

## Features

- Search torrents with various filters (type, category, sort order)
- Get detailed torrent information by ID
- Retrieve Hit & Run activity data
- Get recommended torrents
- Download torrent files directly
- Stateless authentication using a single base64-encoded token
- Optional 2FA support

## Usage

```bash
docker run -p 8080:8080 imdonix/ncore:latest
```

## REST API

The server runs on port 8080 by default.

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/login` | Authenticate and get token |
| `POST` | `/search` | Search torrents |
| `GET` | `/verify` | Verify auth token is validity |
| `GET` | `/torrent/:id` | Get torrent details |
| `GET` | `/torrent/:id/download` | Download torrent file |
| `GET` | `/activity` | Get Hit & Run list |
| `GET` | `/recommended` | Get recommended torrents |
| `POST` | `/logout` | Clear session (stateless) |

### Authentication

**Login**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "your_user", "password": "your_pass", "2factor": "123456"}'
```

```json
{
  "token": "your_token"
}
```

Returns a single `token` for use in subsequent requests.

**Auth Headers**

All authenticated requests must include the token in one of the following headers:
- `X-Ncore-Auth: <token>`

### Search Request

```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/json" \
  -H "X-Ncore-Auth: your_token_here" \
  -d '{
    "pattern": "lord of the rings",
    "type": "all_own",
    "where": "name",
    "sort_by": "seeders",
    "sort_order": "DESC",
    "page": 1
  }'
```

```json
{
  "Torrents": [
    {
      "ID": "3815771",
      "Title": "The.Lord.of.the.Rings.The.Rings.of.Power.S02.AMZN.WEBRiP.AAC2.0.x264.HuN.EnG-B9R",
      "Key": "55f4dac4c72e5e2a28cecb3929f41091",
      "Size": {},
      "Type": "xvidser_hun",
      "Date": "2024-10-03T10:53:52Z",
      "Seeders": 1012,
      "Leechers": 130,
      "Download": "https://ncore.pro/torrents.php?action=download&id=3815771&key=55f4dac4c72e5e2a28cecb3929f41091",
      "URL": "https://ncore.pro/torrents.php?action=details&id=3815771",
      "Extra": null
    },
    ...
  ]
}
```

### Download Torrent

```bash
curl -X GET "http://localhost:8080/torrent/123456/download" \
  -H "X-Ncore-Auth: your_token_here" \
  -o torrent.torrent
```

## Library Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/imdonix/ncore-go/pkg/ncore"
)

func main() {
    // 1. Login to get a token
    client, _ := ncore.NewClient(15*time.Second)
    token, err := client.Login("username", "password", "")
    if err != nil {
        panic(err)
    }
    fmt.Println("Logged in, token:", token)

    // 2. Or initialize a client directly from a saved token
    // client, _ = ncore.NewClientFromToken(15*time.Second, token)

    result, _ := client.Search("movie name", ncore.TypeAllOwn, ncore.WhereName, ncore.SortSeeders, ncore.SeqDesc, 1)
    for _, t := range result.Torrents {
        fmt.Printf("%s - %s (%d seeders)\n", t.Name, t.Size, t.Seeders)
    }
}
```

## Build

```bash
make build    # Build binary to bin/ncore
make docker   # Build Docker image
```
