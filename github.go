// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/wabarc/helper"
	"github.com/wabarc/logger"
)

type github struct {
	client *http.Client
}

func newGitHub() *github {
	return &github{
		client: &http.Client{},
	}
}

func (gh *github) request(ctx context.Context, str string) (b []byte, err error) {
	endpoint := "https://api.github.com/search/issues?per_page=1&sort=created&order=desc&q="
	if repo := os.Getenv("PLAYBACK_GITHUB_REPO"); repo != "" {
		str += "+repo:" + repo
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+str+"+archived", nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("PLAYBACK_GITHUB_PAT"); token != "" {
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

func (gh *github) extract(ctx context.Context, input *url.URL, scope string) (string, error) {
	var re string
	if scope == "ipfs" {
		re = `(?i)https?:\/\/ipfs\.io\/ipfs\/\w{46}`
	} else {
		re = `(?i)https?:\/\/telegra\.ph\/.+?\-\d{2}\-\d{2}`
	}

	data, err := gh.request(ctx, input.String())
	if err != nil {
		logger.Error("[playback] error: %v", err)
		return "", err
	}
	dst := matchLink(re, parseIssue(data))

	if dst == "" {
		return "", errNotFound
	}

	return dst, nil
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
	return ""
}

// nolint:deadcode,unused
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
