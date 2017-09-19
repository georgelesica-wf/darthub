package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var cutRepoURLs []string = make([]string, 0)

func usage() {
	println("darthub <user> <write_file>")
}

func main() {
	flag.Parse()
	userName := flag.Arg(0)
	fileFlag := flag.Arg(1)

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

	if fileFlag != "" {
		for _, result := range results {
		cutRepoURL := strings.Split(result.URL, "github.com/")[1]
		cutRepoURLs = append(cutRepoURLs, cutRepoURL)
	}

		file, err := os.Create("dart_medic_repo_list.py")
		if err == nil {
			data := json.NewEncoder(file)
			data.SetIndent("", "  ")
			err = data.Encode(cutRepoURLs)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}

	for _, result := range results {
		println(result.URL)
	}
}
