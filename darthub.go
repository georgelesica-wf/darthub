package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var cutRepoURLs []string = make([]string, 0)

func usage() {
	println("darthub <user> <write_file>")
}

func main() {
	pythonOutput := flag.Bool("python", false, "Output formatted as a Python list")
	trimURL := flag.Bool("trim", false, "Trim the repo host name from the output")

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
		User:      userName,
		Filename:  "pubspec",
		Extension: "yaml",
		PerPage:   RESULTS_PER_PAGE,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if *pythonOutput {
		println("[")
	}

	for _, result := range results {
		if *pythonOutput {
			print("  \"")
		}

		if *trimURL {
			print(result.TrimmedURL)
		} else {
			print(result.URL)
		}

		if *pythonOutput {
			print("\",")
		}

		println()
	}

	if *pythonOutput {
		println("]")
	}
}
