// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/wabarc/helper"
	"github.com/wabarc/logger"
)

var errGHNotFound = errors.New("Not found")

type github struct {
	client *http.Client
}

func newGitHub() *github {
	return &github{
		client: &http.Client{},
	}
}

func (gh *github) request(str string) (b []byte, err error) {
	endpoint := "https://api.github.com/search/issues?per_page=1&q="

	req, err := http.NewRequest("GET", endpoint+url.QueryEscape(str+" origin archived"), nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	token := os.Getenv("PLAYBACK_GITHUB_PAT")
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := gh.client.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return body, nil
}

func (gh *github) extract(links []string, scope string) (map[string]string, error) {
	collects := collects(links)
	results := make(map[string]string)

	var re string
	var mu sync.Mutex
	var wg sync.WaitGroup

	if scope == "ipfs" {
		re = `(?i)https?:\/\/ipfs\.io\/ipfs\/\w{46}`
	} else {
		re = `(?i)https?:\/\/telegra\.ph\/.+?\-\d{2}\-\d{2}`
	}

	// nolint:staticcheck
	var strip = func(link string) string {
		if u, err := url.Parse(link); err == nil {
			u.Scheme = ""
			link = u.String()
			link = strings.TrimLeft(link, "//")
			link = strings.TrimLeft(link, "www.")
			link = strings.TrimLeft(link, "wap.")
			link = strings.TrimLeft(link, "m.")
		}
		return link
	}

	for _, link := range collects {
		wg.Add(1)
		go func(link string) {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			data, err := gh.request(link)
			if err != nil {
				logger.Error("[playback] error: %v", err)
				results[link] = "Unknow error"
				return
			}
			results[link] = matchLink(re, parseIssue(data))
		}(strip(link))
	}
	wg.Wait()

	if len(results) == 0 {
		return results, errGHNotFound
	}

	return results, nil
}

func matchLink(regex, str string) string {
	var re = regexp.MustCompile(regex)
	for _, match := range re.FindAllString(str, -1) {
		uri, err := url.Parse(match)
		if err != nil {
			continue
		}
		return uri.String()
	}
	return "Not Found"
}

func collects(links []string) map[string]string {
	collects := make(map[string]string)
	for _, link := range links {
		if !helper.IsURL(link) {
			logger.Info("[playback]" + link + " is invalid url.")
			continue
		}
		collects[link] = link
	}
	return collects
}

func parseIssue(data []byte) string {
	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		logger.Debug("[playback] unmarshal json failed: %v", err)
		return ""
	}
	items, ok := dat["items"].([]interface{})
	if !ok {
		logger.Debug("[playback] parse items failed: %v", items)
		return ""
	}
	if len(items) == 0 {
		logger.Debug("[playback] items not found")
		return ""
	}
	item, ok := items[0].(map[string]interface{})
	if !ok {
		logger.Debug("[playback] parse item field failed: %v", item)
		return ""
	}
	body, ok := item["body"].(string)
	if !ok {
		logger.Debug("[playback] parse body field failed: %v", body)
		return ""
	}

	return body
}
