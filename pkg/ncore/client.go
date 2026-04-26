package ncore

import (
	"slices"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	loggedIn   bool
}

func NewClient(timeout time.Duration, cookies map[string]string) (*Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &Client{
		httpClient: &http.Client{
			Timeout: timeout,
			Jar:     jar,
		},
	}

	if cookies != nil {
		u, _ := url.Parse(URLIndex)
		var httpCookies []*http.Cookie
		for name, value := range cookies {
			found := slices.Contains(AllowedCookies, name)
			if found {
				httpCookies = append(httpCookies, &http.Cookie{
					Name:   name,
					Value:  value,
					Domain: URLCookieDomain,
				})
			}
		}
		jar.SetCookies(u, httpCookies)
		if client.checkLoggedIn() {
			client.loggedIn = true
		}
	}

	return client, nil
}

func (c *Client) checkLoggedIn() bool {
	resp, err := c.httpClient.Get(URLIndex)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if strings.Contains(resp.Request.URL.String(), "login.php") {
		return false
	}

	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "<title>nCore</title>") {
		return false
	}

	return true
}

func (c *Client) Login(username, password, twoFactorCode string) (map[string]string, error) {
	if c.loggedIn && c.checkLoggedIn() {
		return c.getCookies(), nil
	}

	c.httpClient.Jar, _ = cookiejar.New(nil)
	c.loggedIn = false

	data := url.Values{}
	data.Set("nev", username)
	data.Set("pass", password)
	data.Set("set_lang", "hu")
	data.Set("submitted", "1")
	data.Set("ne_leptessen_ki", "1")
	if twoFactorCode != "" {
		data.Set("2factor", twoFactorCode)
	}

	resp, err := c.httpClient.PostForm(URLLogin, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.Request.URL.String() != URLIndex || strings.Contains(string(body), "<title>nCore</title>") {
		c.Logout()
		errorMsg := fmt.Sprintf("check credentials for user: '%s'", username)
		if twoFactorCode != "" {
			errorMsg += ". Invalid 2FA code or wait 5 minutes between login attempts."
		}
		return nil, fmt.Errorf("%w: %s", ErrLoginFailed, errorMsg)
	}

	c.loggedIn = true
	return c.getCookies(), nil
}

func (c *Client) getCookies() map[string]string {
	u, _ := url.Parse(URLIndex)
	cookies := make(map[string]string)
	for _, cookie := range c.httpClient.Jar.Cookies(u) {
		for _, allowed := range AllowedCookies {
			if cookie.Name == allowed {
				cookies[cookie.Name] = cookie.Value
			}
		}
	}
	return cookies
}

func (c *Client) Search(pattern string, tType SearchParamType, where SearchParamWhere, sortBy ParamSort, sortOrder ParamSeq, page int) (*SearchResult, error) {
	if !c.loggedIn {
		return nil, ErrNotLoggedIn
	}

	searchURL := fmt.Sprintf(URLDownloadPattern, page, string(tType), string(sortBy), string(sortOrder), url.QueryEscape(pattern), string(where))
	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	torrents, err := parseTorrentsPage(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParserError, err)
	}

	numOfPages := getNumOfPages(bodyStr)
	return &SearchResult{Torrents: torrents, NumOfPages: numOfPages}, nil
}

func (c *Client) GetTorrent(id string) (*Torrent, error) {
	if !c.loggedIn {
		return nil, ErrNotLoggedIn
	}

	detailURL := fmt.Sprintf(URLDetailPattern, id)
	resp, err := c.httpClient.Get(detailURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return parseTorrentDetail(string(body), id)
}

func (c *Client) GetByActivity() ([]*Torrent, error) {
	if !c.loggedIn {
		return nil, ErrNotLoggedIn
	}

	resp, err := c.httpClient.Get(URLActivity)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	params := parseActivity(string(body))

	var torrents []*Torrent
	for _, p := range params {
		if len(p) < 8 {
			continue
		}
		id := p[0]
		uploaded, _ := NewSize(p[4])
		downloaded, _ := NewSize(p[5])
		rate, _ := strconv.ParseFloat(p[7], 64)

		t, err := c.GetTorrent(id)
		if err == nil {
			t.Extra = map[string]any{
				"start":      p[1],
				"updated":    p[2],
				"status":     p[3],
				"uploaded":   uploaded,
				"downloaded": downloaded,
				"remaining":  p[6],
				"rate":       rate,
			}
			torrents = append(torrents, t)
		}
	}
	return torrents, nil
}

func (c *Client) GetRecommended(tType SearchParamType) ([]*Torrent, error) {
	if !c.loggedIn {
		return nil, ErrNotLoggedIn
	}

	resp, err := c.httpClient.Get(URLRecommended)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ids := parseRecommendedIds(string(body))

	var torrents []*Torrent
	for _, id := range ids {
		t, err := c.GetTorrent(id)
		if err == nil {
			if tType == "" || tType == TypeAllOwn || t.Type == tType {
				torrents = append(torrents, t)
			}
		}
	}
	return torrents, nil
}

func (c *Client) Download(torrent *Torrent) (io.ReadCloser, string, error) {
	if !c.loggedIn {
		return nil, "", ErrNotLoggedIn
	}

	filename, downloadURL := torrent.PrepareDownload("")

	resp, err := c.httpClient.Get(downloadURL)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrConnectionError, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", fmt.Errorf("%w: unexpected status code: %d", ErrDownloadFailed, resp.StatusCode)
	}

	return resp.Body, filename, nil
}

func (c *Client) Logout() {
	c.httpClient.Jar, _ = cookiejar.New(nil)
	c.loggedIn = false
}
