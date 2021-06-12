// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package playback

import (
	"testing"
)

func TestPlayback(t *testing.T) {
	tests := []struct {
		name string
		urls []string
		got  int
	}{
		{
			name: "Without URLs",
			urls: []string{},
			got:  0,
		},
		{
			name: "Has one invalid URL",
			urls: []string{"foo bar", "https://example.com/"},
			got:  1,
		},
		{
			name: "URLs full matches",
			urls: []string{"https://example.com/", "https://example.org/"},
			got:  2,
		},
	}

	for _, test := range tests {
		t.Run("IA_"+test.name, func(t *testing.T) {
			var pb Playback = &Handle{URLs: test.urls}
			got := pb.IA()
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
		t.Run("IS_"+test.name, func(t *testing.T) {
			var pb Playback = &Handle{URLs: test.urls}
			got := pb.IS()
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
		t.Run("PH_"+test.name, func(t *testing.T) {
			var pb Playback = &Handle{URLs: test.urls}
			got := pb.PH()
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
		t.Run("IP_"+test.name, func(t *testing.T) {
			var pb Playback = &Handle{URLs: test.urls}
			got := pb.IP()
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
		t.Run("TT_"+test.name, func(t *testing.T) {
			var pb Playback = &Handle{URLs: test.urls}
			got := pb.TT()
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
	}
}

func TestExtractIPFSLink(t *testing.T) {
	var got map[string]string

	tests := []struct {
		name string
		urls []string
		got  int
	}{
		{
			name: "Without URLs",
			urls: []string{},
			got:  0,
		},
		{
			name: "Has one invalid URL",
			urls: []string{"foo bar", "https://example.com/"},
			got:  1,
		},
		{
			name: "URLs full matches",
			urls: []string{"https://example.com/", "https://example.org/"},
			got:  2,
		},
	}

	gh := newGitHub()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ = gh.extract(test.urls, "ipfs")
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
	}
}

func TestExtractTelegraphLink(t *testing.T) {
	var got map[string]string

	tests := []struct {
		name string
		urls []string
		got  int
	}{
		{
			name: "Without URLs",
			urls: []string{},
			got:  0,
		},
		{
			name: "Has one invalid URL",
			urls: []string{"foo bar", "https://example.com/"},
			got:  1,
		},
		{
			name: "URLs full matches",
			urls: []string{"https://example.com/", "https://example.org/"},
			got:  2,
		},
	}

	gh := newGitHub()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ = gh.extract(test.urls, "telegraph")
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
	}
}

func TestGoogleCache(t *testing.T) {
	var got map[string]string

	tests := []struct {
		name string
		urls []string
		got  int
	}{
		{
			name: "Without URLs",
			urls: []string{},
			got:  0,
		},
		{
			name: "Has one invalid URL",
			urls: []string{"foo bar", "https://example.com/"},
			got:  1,
		},
		{
			name: "URLs full matches",
			urls: []string{"https://example.com/", "https://example.org/"},
			got:  2,
		},
	}

	g := google()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ = g.cache(test.urls)
			if len(got) != test.got {
				t.Errorf("got = %d; want %d", len(got), test.got)
			}
			for orig, dest := range got {
				if testing.Verbose() {
					t.Log(orig, "=>", dest)
				}
			}
		})
	}
}
