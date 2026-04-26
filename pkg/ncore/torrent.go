package ncore

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Torrent struct {
	ID       string
	Title    string
	Key      string
	Size     Size
	Type     SearchParamType
	Date     time.Time
	Seeders  int
	Leechers int
	Download string
	URL      string
	Extra    map[string]any
}

func (t *Torrent) PrepareDownload(path string) (string, string) {
	filename := strings.ReplaceAll(t.Title, " ", "_") + ".torrent"
	filePath := filepath.Join(path, filename)
	return filePath, t.Download
}

func (t *Torrent) String() string {
	return fmt.Sprintf("<Torrent %s>", t.ID)
}

type SearchResult struct {
	Torrents   []*Torrent
	NumOfPages int
}
