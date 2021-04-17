package main

import (
	"flag"
	"fmt"
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
		stub map[string]string
	}

	var pb playback.Playback = &playback.Handle{URLs: args}

	results := []collects{
		{slot: "Internet Archive", stub: pb.IA()},
		{slot: "archive.today", stub: pb.IS()},
		{slot: "IPFS", stub: pb.IP()},
		{slot: "Telegraph", stub: pb.PH()},
		{slot: "Time Travel", stub: pb.TT()},
	}
	for _, collect := range results {
		fmt.Printf("[%s]\n", collect.slot)
		for orig, dest := range collect.stub {
			fmt.Println(orig, "=>", dest)
		}
		fmt.Printf("\n")
	}
}
