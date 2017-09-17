package main

import (
	"github.com/google/go-github/github"
	"context"
	"golang.org/x/oauth2"
	"os"
	"flag"
	"fmt"
)

func usage() {
	println("darthub <user>")
}

func main() {
	flag.Parse()
	userName := flag.Arg(0)

	if userName == "" {
		usage()
		println("User name parameter is required")
		return
	}

	token, present := os.LookupEnv("DARTHUB_TOKEN")
	if !present {
		println("DARTHUB_TOKEN is not set")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	searcher := GithubSearcher{client: client}

	results, err := searcher.Search(SearchParams{
		User: userName,
		Filename: "pubspec",
		Extension: "yaml",
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		println(result.URL)
	}
}
