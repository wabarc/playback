package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/wabarc/playback"
)

func main() {
	var (
		version bool
	)

	const versionHelp = "Show version"

	flag.BoolVar(&version, "version", false, versionHelp)
	flag.BoolVar(&version, "v", false, versionHelp)
	flag.Parse()

	if version {
		fmt.Println(playback.Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		e := os.Args[0]
		fmt.Printf("  %s url [url]\n\n", e)
		fmt.Printf("example:\n  %s https://example.com https://example.org\n\n", e)
		os.Exit(1)
	}

	type collects struct {
		slot string
		stub playback.Playbacker
	}

	var wrap = func(input *url.URL) []collects {
		return []collects{
			{slot: "Internet Archive", stub: playback.IA{URL: input}},
			{slot: "archive.today", stub: playback.IS{URL: input}},
			{slot: "Ghostarchive", stub: playback.GA{URL: input}},
			{slot: "IPFS", stub: playback.IP{URL: input}},
			{slot: "Telegraph", stub: playback.PH{URL: input}},
			{slot: "Time Travel", stub: playback.TT{URL: input}},
		}
	}
	for _, arg := range args {
		input, err := url.Parse(arg)
		if err != nil {
			fmt.Println(arg, "=>", fmt.Sprint(err))
			continue
		}
		for _, collect := range wrap(input) {
			fmt.Printf("[%s]\n", collect.slot)
			dest := playback.Playback(context.TODO(), collect.stub)
			fmt.Println(arg, "=>", dest)
			fmt.Printf("\n")
		}
	}
}
