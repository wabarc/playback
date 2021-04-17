// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

/*
Playback is a toolkit for search webpages archived
to Internet Archive, archive.today, IPFS and beyond.
*/

package playback // import "github.com/wabarc/playback"

import (
	"github.com/wabarc/archive.is"
	"github.com/wabarc/archive.org"
	"github.com/wabarc/logger"
	"github.com/wabarc/memento"
)

// Archives represents result from the time capsules.
type Archives map[string]string

// Playback is interface of the playback,
// methods returns `Archives`.
type Playback interface {
	IA() Archives // Internet Archive
	IS() Archives // archive.today
	PH() Archives // Telegra.ph
	IP() Archives // IPFS
	TT() Archives // Time Travel, http://timetravel.mementoweb.org
}

// Handle represents a playback handle.
type Handle struct {
	URLs []string
}

func (h *Handle) IA() Archives {
	wbrc := &ia.Archiver{}
	uris, err := wbrc.Playback(h.URLs)
	if err != nil {
		logger.Error("Playback %v from Internet Archive failed, %v", h.URLs, err)
	}

	return uris
}

func (h *Handle) IS() Archives {
	wbrc := &is.Archiver{}
	uris, err := wbrc.Playback(h.URLs)
	if err != nil {
		logger.Error("Playback %v from archive.today failed, %v", h.URLs, err)
	}

	return uris
}

func (h *Handle) IP() Archives {
	uris, err := newGitHub().extract(h.URLs, "ipfs")
	if err != nil {
		logger.Error("Playback %v from IPFS failed, %v", h.URLs, err)
	}

	return uris
}

func (h *Handle) PH() Archives {
	uris, err := newGitHub().extract(h.URLs, "telegraph")
	if err != nil {
		logger.Error("Playback %v from IPFS failed, %v", h.URLs, err)
	}

	return uris
}

func (h *Handle) TT() Archives {
	wbrc := &memento.Memento{}
	uris, err := wbrc.Mementos(h.URLs)
	if err != nil {
		logger.Error("Playback %v from Time Travel failed, %v", h.URLs, err)
	}

	return uris
}
