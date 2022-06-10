// Copyright 2022 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback // import "github.com/wabarc/playback"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/wabarc/logger"
)

var defaultIndexing = "capsules"

type meili struct {
	disabled bool
	endpoint string
	apikey   string

	client *http.Client
}

func newMeili() *meili {
	endpoint := os.Getenv("PLAYBACK_MEILI_ENDPOINT")
	indexing := os.Getenv("PLAYBACK_MEILI_INDEXING")
	apikey := os.Getenv("PLAYBACK_MEILI_APIKEY")
	if indexing == "" {
		indexing = defaultIndexing
	}
	disabled := endpoint == ""
	endpoint = fmt.Sprintf("%s/indexes/%s/search", endpoint, indexing)
	return &meili{
		disabled: disabled,
		endpoint: endpoint,
		apikey:   apikey,
		client:   &http.Client{},
	}
}

func (m *meili) request(ctx context.Context, str string) (b []byte, err error) {
	if m.disabled {
		return nil, errors.New(`meilisearch disabled`)
	}

	// doc: https://docs.meilisearch.com/reference/api/search.html
	params := struct {
		Query                string   `json:"q"`
		Limit                int      `json:"limit"`
		Sort                 []string `json:"sort"`
		Matches              bool     `json:"matches"`
		AttributesToRetrieve []string `json:"attributesToRetrieve"`
	}{
		Query:                fmt.Sprintf(`"%s"`, str),
		Limit:                1,
		Sort:                 []string{"id:desc"},
		Matches:              true,
		AttributesToRetrieve: []string{"ip", "ph"},
	}
	buf, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(buf)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.endpoint, payload)
	req.Header.Add("Content-Type", "application/json")
	if m.apikey != "" {
		req.Header.Add("Authorization", "Bearer "+m.apikey)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return body, nil
}

func (m *meili) extract(ctx context.Context, input *url.URL, scope string) (dst string, err error) {
	type response struct {
		Hits []struct {
			ID     string `json:"id"`
			Source string `json:"source"`
			IA     string `json:"ia"`
			IS     string `json:"is"`
			IP     string `json:"ip"`
			PH     string `json:"ph"`
		} `json:"hits"`
	}

	// Remove scheme
	str := strings.TrimLeft(input.String(), input.Scheme)
	resp, err := m.request(ctx, str)
	if err != nil {
		logger.Error("playback from meilisearch failed: %v", err)
		return "", err
	}

	var data response
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return "", err
	}
	if len(data.Hits) == 0 {
		return "", errNotFound
	}

	switch scope {
	case "ip", "ipfs":
		return data.Hits[0].IP, nil
	case "ph", "telegraph":
		return data.Hits[0].PH, nil
	}

	return "", errNotFound
}
