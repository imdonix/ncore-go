package ncore

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

var (
	typePattern        = regexp.MustCompile(`<a href=".*\/torrents\.php\?tipus=(.*?)"><img src=".*" class="categ_link" alt=".*" title=".*">`)
	idAndNamePattern   = regexp.MustCompile(`<a href=".*?" onclick="torrent\(([0-9]+)\); return false;" title="(.*?)">`)
	dateAndTimePattern = regexp.MustCompile(`<div class="box_feltoltve2">(.*?)<br>(.*?)</div>`)
	sizePattern        = regexp.MustCompile(`<div class="box_meret2">(.*?)</div>`)
	seedersPattern     = regexp.MustCompile(`<div class="box_s2"><a class="torrent" href=".*">([0-9]+)</a></div>`)
	leechersPattern    = regexp.MustCompile(`<div class="box_l2"><a class="torrent" href=".*">([0-9]+)</a></div>`)
	currentPagePattern = regexp.MustCompile(`<span class="active_link"><strong>(\d+).*?</strong></span>`)
	lastPagePattern    = regexp.MustCompile(`<a href="/torrents\.php\?oldal=(\d+)[^>]*><strong>Utolsó</strong></a>`)
	keyPattern         = regexp.MustCompile(`<link rel="alternate" href=".*?\/rss.php\?key=(?P<key>[a-z,0-9]+)" title=".*"`)
	notFoundPattern    = regexp.MustCompile(`<div class="lista_mini_error">Nincs találat!</div>`)

	detailTypePattern  = regexp.MustCompile(`<div class="dd"><a title=".*?" href=".*?torrents.php\?csoport_listazas=(?P<category>.*?)">.*?</a>.*?<a title=".*?" href=".*?torrents.php\?tipus=(?P<type>.*?)">.*?</a></div>`)
	detailDatePattern  = regexp.MustCompile(`<div class="dd">(?P<date>[0-9]{4}\-[0-9]{2}\-[0-9]{2}\ [0-9]{2}\:[0-9]{2}\:[0-9]{2})</div>`)
	detailTitlePattern = regexp.MustCompile(`<div class="torrent_reszletek_cim">(?P<title>.*?)</div>`)
	detailSizePattern  = regexp.MustCompile(`<div class="dd">(?P<size>[0-9,.]+\ [K,M,G,T]{1}iB)\ \(.*?\)</div>`)
	detailPeersPattern = regexp.MustCompile(`(?s)div class="dt">Seederek:</div>.*?<div class="dd"><a onclick=".*?">(?P<seed>[0-9]+)</a></div>.*?<div class="dt">Leecherek:</div>.*?<div class="dd"><a onclick=".*?">(?P<leech>[0-9]+)</a></div>`)

	recommendedIdPattern = regexp.MustCompile(`<a href=".*?torrents.php\?action=details\&id=(.*?)" target=".*?"><img src=".*?" width=".*?" height=".*?" border=".*?" title=".*?"\/><\/a>`)

	activityPatterns = []*regexp.Regexp{
		regexp.MustCompile(`onclick="torrent\((.*?)\);`),
		regexp.MustCompile(`<div class="hnr_tstart">(.*?)<\/div>`),
		regexp.MustCompile(`<div class="hnr_tlastactive">(.*?)<\/div>`),
		regexp.MustCompile(`<div class="hnr_tseed"><span class=".*?">(.*?)<\/span><\/div>`),
		regexp.MustCompile(`<div class="hnr_tup">(.*?)<\/div>`),
		regexp.MustCompile(`<div class="hnr_tdown">(.*?)<\/div>`),
		regexp.MustCompile(`<div class="hnr_ttimespent"><span class=".*?">(.*?)<\/span><\/div>`),
		regexp.MustCompile(`<div class="hnr_tratio"><span class=".*?">(.*?)<\/span><\/div>`),
	}
)

func getKey(data string) (string, error) {
	match := keyPattern.FindStringSubmatch(data)
	if len(match) > 0 {
		for i, name := range keyPattern.SubexpNames() {
			if name == "key" {
				return match[i], nil
			}
		}
	}
	return "", fmt.Errorf("error while read user key")
}

func parseTorrentsPage(data string) ([]*Torrent, error) {
	types := typePattern.FindAllStringSubmatch(data, -1)
	idsAndNames := idAndNamePattern.FindAllStringSubmatch(data, -1)
	datesAndTimes := dateAndTimePattern.FindAllStringSubmatch(data, -1)
	sizes := sizePattern.FindAllStringSubmatch(data, -1)
	seeds := seedersPattern.FindAllStringSubmatch(data, -1)
	leeches := leechersPattern.FindAllStringSubmatch(data, -1)

	if len(types) == 0 {
		if notFoundPattern.MatchString(data) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while parse download items")
	}

	if len(types) != len(idsAndNames) || len(types) != len(datesAndTimes) || len(types) != len(sizes) || len(types) != len(seeds) || len(types) != len(leeches) {
		return nil, fmt.Errorf("mismatch in parsed items lengths")
	}

	key, err := getKey(data)
	if err != nil {
		return nil, err
	}

	var torrents []*Torrent
	for i := range types {
		id := idsAndNames[i][1]
		title := idsAndNames[i][2]
		dateStr := datesAndTimes[i][1]
		timeStr := datesAndTimes[i][2]
		dateTime, _ := time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
		size, _ := NewSize(sizes[i][1])
		seed, _ := strconv.Atoi(seeds[i][1])
		leech, _ := strconv.Atoi(leeches[i][1])

		torrents = append(torrents, &Torrent{
			ID:       id,
			Title:    title,
			Key:      key,
			Date:     dateTime,
			Size:     size,
			Type:     SearchParamType(types[i][1]),
			Seeders:  seed,
			Leechers: leech,
			Download: fmt.Sprintf(URLDownloadLink, id, key),
			URL:      fmt.Sprintf(URLDetailPattern, id),
		})
	}
	return torrents, nil
}

func getNumOfPages(data string) int {
	currentMatch := currentPagePattern.FindStringSubmatch(data)
	lastMatch := lastPagePattern.FindStringSubmatch(data)

	numOfPages := 0
	if len(currentMatch) > 1 {
		currentItems, _ := strconv.Atoi(currentMatch[1])
		numOfPages = int(math.Ceil(float64(currentItems) / 25.0))
	}
	if len(lastMatch) > 1 {
		lastPage, _ := strconv.Atoi(lastMatch[1])
		if lastPage > numOfPages {
			numOfPages = lastPage
		}
	}
	return numOfPages
}

func parseTorrentDetail(data string, id string) (*Torrent, error) {
	tTypeMatch := detailTypePattern.FindStringSubmatch(data)
	if tTypeMatch == nil {
		return nil, fmt.Errorf("type pattern not found")
	}

	category := ""
	tType := ""
	for i, name := range detailTypePattern.SubexpNames() {
		switch name {
			case "category":
				category = tTypeMatch[i]
			case "type":
				tType = tTypeMatch[i]
		}
	}

	finalType, ok := detailedParamMap[category+"_"+tType]
	if !ok {
		// Fallback or error
		finalType = SearchParamType(tType)
	}

	dateMatch := detailDatePattern.FindStringSubmatch(data)
	if dateMatch == nil {
		return nil, fmt.Errorf("date pattern not found")
	}
	dateTime, _ := time.Parse("2006-01-02 15:04:05", dateMatch[1])

	titleMatch := detailTitlePattern.FindStringSubmatch(data)
	if titleMatch == nil {
		return nil, fmt.Errorf("title pattern not found")
	}
	title := titleMatch[1]

	key, err := getKey(data)
	if err != nil {
		return nil, err
	}

	sizeMatch := detailSizePattern.FindStringSubmatch(data)
	if sizeMatch == nil {
		return nil, fmt.Errorf("size pattern not found")
	}
	size, _ := NewSize(sizeMatch[1])

	peersMatch := detailPeersPattern.FindStringSubmatch(data)
	seed := 0
	leech := 0
	if peersMatch != nil {
		for i, name := range detailPeersPattern.SubexpNames() {
			switch name {
				case "seed":
					seed, _ = strconv.Atoi(peersMatch[i])
				case "leech":
					leech, _ = strconv.Atoi(peersMatch[i])
			}
		}
	}

	return &Torrent{
		ID:       id,
		Title:    title,
		Key:      key,
		Date:     dateTime,
		Size:     size,
		Type:     finalType,
		Seeders:  seed,
		Leechers: leech,
		Download: fmt.Sprintf(URLDownloadLink, id, key),
		URL:      fmt.Sprintf(URLDetailPattern, id),
	}, nil
}

func parseRecommendedIds(data string) []string {
	matches := recommendedIdPattern.FindAllStringSubmatch(data, -1)
	var ids []string
	for _, m := range matches {
		ids = append(ids, m[1])
	}
	return ids
}

func parseActivity(data string) [][]string {
	var results [][]string

	// This is a bit more complex in Python: zip(*out)
	// We need to find all matches for each pattern and then group them by index.

	var allMatches [][]string
	for _, re := range activityPatterns {
		matches := re.FindAllStringSubmatch(data, -1)
		var m []string
		for _, match := range matches {
			m = append(m, match[1])
		}
		allMatches = append(allMatches, m)
	}

	if len(allMatches) == 0 || len(allMatches[0]) == 0 {
		return nil
	}

	numItems := len(allMatches[0])
	for i := 0; i < numItems; i++ {
		item := make([]string, len(allMatches))
		for j := 0; j < len(allMatches); j++ {
			if i < len(allMatches[j]) {
				item[j] = allMatches[j][i]
			}
		}
		results = append(results, item)
	}

	return results
}
