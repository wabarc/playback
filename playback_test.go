// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestPlayback(t *testing.T) {
	t.Parallel()

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
			got := Playback(context.TODO(), test.play)
			if got == "" {
				t.Errorf("playback empty")
			}
		})
	}
}

func TestExtractIPFSLink(t *testing.T) {
	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	gh := newGitHub()
	got, err := gh.extract(context.TODO(), in, "ipfs")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "ipfs.io") {
		t.Log(uri, "=>", got)
		t.Errorf("Unexpect extract ipfs link, got %s does not contains ipfs.io", got)
	}
}

func TestExtractTelegraphLink(t *testing.T) {
	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	gh := newGitHub()
	got, err := gh.extract(context.TODO(), in, "telegraph")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "telegra.ph") {
		t.Log(uri, "=>", got)
		t.Errorf("Unexpect extract telegra.ph link, got %s does not contains telegra.ph", got)
	}
}

func TestGoogleCache(t *testing.T) {
	uri := "https://example.com"
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	got, err := newGoogle().cache(context.TODO(), in)
	if err != nil {
		t.Log(uri, "=>", got)
		t.Fatal(err)
	}
}
