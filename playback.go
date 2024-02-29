// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

/*
Playback is a toolkit for search webpages archived
to Internet Archive, archive.today, IPFS and beyond.
*/

package playback // import "github.com/wabarc/playback"

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/wabarc/archive.is"
	"github.com/wabarc/archive.org"
	"github.com/wabarc/ghostarchive"
	"github.com/wabarc/logger"
	"github.com/wabarc/memento"
)

var errNotFound = fmt.Errorf("Not found")

func init() {
	debug := os.Getenv("DEBUG")
	if debug == "true" || debug == "1" || debug == "on" {
		logger.EnableDebug()
	}
}

// Collect represents a collections from the time capsules.
type Collect struct {
	// Arc string // Archive slot name, see config/config.go
	Dst string // Archived destination URL
	Src string // Source URL
	// Ext string // Extra identifier
}

type IA struct {
	URL *url.URL
}

type IS struct {
	URL *url.URL
}

// IPFS
type IP struct {
	URL *url.URL
}

type PH struct {
	URL *url.URL
}

type GA struct {
	URL *url.URL
}

// Time Travel, http://timetravel.mementoweb.org
type TT struct {
	URL *url.URL
}

// Google Cache
type GC struct {
	URL *url.URL
}

// Playbacker is the interface that wraps the basic Playback method.
//
// Playback playback *url.URL from implementations from the Wayback Machine.
// It returns the result of string from the upstream service.
type Playbacker interface {
	Playback(ctx context.Context) string
}

func (i IA) Playback(ctx context.Context) string {
	arc := &ia.Archiver{}
	dst, err := arc.Playback(ctx, i.URL)
	if err != nil {
		logger.Error("[playback] %s from Internet Archive failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}

	return dst
}

func (i IS) Playback(ctx context.Context) string {
	arc := &is.Archiver{}
	dst, err := arc.Playback(ctx, i.URL)
	if err != nil {
		logger.Error("[playback] %s from archive.today failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}

	return dst
}

func (g GA) Playback(ctx context.Context) string {
	arc := &ga.Archiver{}
	dst, err := arc.Playback(ctx, g.URL)
	if err != nil {
		logger.Error("[playback] %s from Ghostarchive failed: %v", g.URL.String(), err)
		return fmt.Sprint(err)
	}

	return dst
}

func (i IP) Playback(ctx context.Context) string {
	dst, err := newGitHub().extract(ctx, i.URL, "ipfs")
	if err != nil && err != errNotFound {
		logger.Error("[playback] %s from IPFS failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}
	if dst == "" {
		dst, err = newMeili().extract(ctx, i.URL, "ipfs")
	}
	if err != nil {
		return fmt.Sprint(err)
	}

	return dst
}

func (i PH) Playback(ctx context.Context) string {
	dst, err := newGitHub().extract(ctx, i.URL, "telegraph")
	if err != nil && err != errNotFound {
		logger.Error("[playback] %s from Telegra.ph failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}
	if dst == "" {
		dst, err = newMeili().extract(ctx, i.URL, "telegraph")
	}
	if err != nil {
		return fmt.Sprint(err)
	}

	return dst
}

func (i TT) Playback(ctx context.Context) string {
	arc := &memento.Memento{}
	dst, err := arc.Mementos(ctx, i.URL)
	if err != nil {
		logger.Error("[playback] %s from Time Travel failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}

	return dst
}

func (i GC) Playback(ctx context.Context) string {
	dst, err := newGoogle().cache(ctx, i.URL)
	if err != nil {
		logger.Error("[playback] %s from Google Cache failed: %v", i.URL.String(), err)
		return fmt.Sprint(err)
	}

	return dst
}

func Playback(ctx context.Context, p Playbacker) string {
	return p.Playback(ctx)
}
