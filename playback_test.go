// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wabarc/helper"
)

var (
	meiliSearchResp = `{
    "hits": [
        {
            "ip": "https://ipfs.io/ipfs/bafybeibndw52abcyf5i672uw123fh2l6ifhwjacjs56fhjk4udeflejrqi",
            "ph": "https://telegra.ph/Example-01-01",
            "_matchesInfo": {
                "ip": [
                    {
                        "start": 0,
                        "length": 5
                    }
                ],
                "ph": [
                    {
                        "start": 0,
                        "length": 5
                    }
                ]
            }
        }
    ],
    "nbHits": 2,
    "exhaustiveNbHits": false,
    "query": "\"https://example.com\"",
    "limit": 1,
    "offset": 0,
    "processingTimeMs": 14
}`
	githubResp = `{
  "total_count": 1,
  "incomplete_results": false,
  "items": [
    {
      "id": 35802,
      "node_id": "MDU6SXNzdWUzNTgwMg==",
      "number": 132,
      "title": "Line Number Indexes Beyond 20 Not Displayed",
      "user": null,
      "labels": [],
      "state": "open",
      "assignee": null,
      "milestone": null,
      "comments": 15,
      "created_at": "2009-07-12T20:10:41Z",
      "updated_at": "2009-07-19T09:23:43Z",
      "closed_at": null,
      "pull_request": null,
      "body": "https://ipfs.io/ipfs/bafybeibndw52abcyf5i672uw123fh2l6ifhwjacjs56fhjk4udeflejrqi  https://telegra.ph/Example-01-01",
      "score": 1,
      "locked": true,
      "author_association": "COLLABORATOR"
    }
  ]
}`
)

func testServer() (*http.Client, *httptest.Server) {
	httpClient, mux, server := helper.MockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/indexes/%s/search", defaultIndexing):
			_, _ = w.Write([]byte(meiliSearchResp))
		case "/search/issues":
			_, _ = w.Write([]byte(githubResp))
		}
	})

	return httpClient, server
}

func TestPlayback(t *testing.T) {
	t.Parallel()

	_, server := testServer()
	defer server.Close()

	os.Setenv("PLAYBACK_MEILI_ENDPOINT", server.URL)

	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		play Playbacker
	}{
		{
			name: "Internet Archive",
			play: IA{URL: in},
		},
		{
			name: "archive.today",
			play: IS{URL: in},
		},
		{
			name: "IPFS",
			play: IP{URL: in},
		},
		{
			name: "Telegra.ph",
			play: PH{URL: in},
		},
		{
			name: "Time Travel",
			play: TT{URL: in},
		},
		{
			name: "Google Cache",
			play: GC{URL: in},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			got := Playback(ctx, test.play)
			if got == "" {
				t.Errorf("playback empty")
			}
		})
	}
}

func TestExtractIPFSLink(t *testing.T) {
	client, server := testServer()
	defer server.Close()

	os.Setenv("PLAYBACK_MEILI_ENDPOINT", server.URL)

	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	gh := newGitHub()
	gh.client = client
	got, err := gh.extract(ctx, in, "ipfs")
	if err != nil && err != errNotFound {
		t.Fatal(err)
	}

	if !strings.Contains(got, "ipfs.io") {
		t.Log(uri, "=>", got)
		t.Errorf("Unexpected extract ipfs link, got %s does not contains ipfs.io", got)
	}
}

func TestExtractTelegraphLink(t *testing.T) {
	client, server := testServer()
	defer server.Close()

	os.Setenv("PLAYBACK_MEILI_ENDPOINT", server.URL)

	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	gh := newGitHub()
	gh.client = client
	got, err := gh.extract(ctx, in, "telegraph")
	if err != nil && err != errNotFound {
		t.Fatal(err)
	}

	if !strings.Contains(got, "telegra.ph") {
		t.Errorf("Unexpected extract telegra.ph link, got %s does not contains telegra.ph", got)
	}
}

func TestGoogleCache(t *testing.T) {
	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	got, err := newGoogle().cache(ctx, in)
	if err != nil {
		t.Log(uri, "=>", got)
		t.Fatal(err)
	}
}

func TestMeilisearch(t *testing.T) {
	_, server := testServer()
	defer server.Close()

	os.Setenv("PLAYBACK_MEILI_ENDPOINT", server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	m := newMeili()
	got, err := m.extract(ctx, in, "telegraph")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got, "telegra.ph") {
		t.Log(uri, "=>", got)
		t.Errorf("Unexpected extract telegra.ph link, got %s does not contains telegra.ph", got)
	}
}
