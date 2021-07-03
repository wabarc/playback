// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback // import "github.com/wabarc/playback"

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/wabarc/logger"
)

var (
	errGCNotFound = errors.New("Not found")

	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"
)

type google struct {
	client *http.Client
}

func newGoogle() *google {
	return &google{
		client: &http.Client{Timeout: time.Minute, CheckRedirect: noRedirect},
	}
}

func (g *google) cache(ctx context.Context, input *url.URL) (string, error) {
	dst, err := g.request(ctx, input.String())
	if err != nil {
		logger.Error("[playback] from google cache error: %v", err)
		return "", err
	}

	if dst == nil {
		return "", errGCNotFound
	}

	return dst.String(), nil
}

func (g *google) request(ctx context.Context, uri string) (dest *url.URL, err error) {
	endpoint := "https://webcache.googleusercontent.com/search?q=cache:"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+uri, nil)

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
