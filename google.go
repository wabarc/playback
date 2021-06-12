// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback // import "github.com/wabarc/playback"

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/wabarc/logger"
)

var (
	errGCNotFound = errors.New("Not found")

	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"
)

type Google struct {
	client *http.Client
}

func google() *Google {
	return &Google{
		client: &http.Client{Timeout: time.Minute, CheckRedirect: noRedirect},
	}
}

func (g *Google) cache(links []string) (map[string]string, error) {
	collects := collects(links)
	results := make(map[string]string)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, link := range collects {
		wg.Add(1)
		go func(link string) {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			dest, err := g.request(link)
			if err != nil {
				logger.Error("[playback] from google cache error: %v", err)
				results[link] = err.Error()
				return
			}
			results[link] = dest.String()
		}(link)
	}
	wg.Wait()

	if len(results) == 0 {
		return results, errGCNotFound
	}

	return results, nil
}

func (g *Google) request(uri string) (dest *url.URL, err error) {
	endpoint := "https://webcache.googleusercontent.com/search?q=cache:"
	req, err := http.NewRequest("GET", endpoint+uri, nil)

	req.Header.Add("User-Agent", userAgent)

	resp, err := g.client.Do(req)
	if err != nil {
		logger.Error("[playback] google cache error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Debug("[playback] google cache status code: %d", resp.StatusCode)
		return nil, fmt.Errorf(resp.Status)
	}

	return resp.Request.URL, nil
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
