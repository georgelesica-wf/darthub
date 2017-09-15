package main

import (
	"github.com/google/go-github/github"
	"context"
	"golang.org/x/oauth2"
	"os"
	"flag"
	"time"
)

const RESULTS_PER_PAGE = 100

var repoURLs []string = make([]string, 0)

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if needle == candidate {
			return true
		}
	}
	return false
}

func fetchResults(client *github.Client, ctx context.Context, userName string, page int) {
	opts := &github.SearchOptions{ListOptions: github.ListOptions{Page: page, PerPage: RESULTS_PER_PAGE}}

	// list all repositories for the authenticated userName
	result, _, err := client.Search.Code(ctx, "user:" +userName+ " filename:pubspec extension:yaml", opts)
	if err != nil {
		panic(err)
	}

	for _, item := range result.CodeResults {
		repoURL := item.Repository.GetHTMLURL()

		if contains(repoURLs, repoURL) {
			continue
		}

		repoURLs = append(repoURLs, repoURL)
		println(repoURL)
	}

	// We are finished when we get back fewer than the max results.
	if len(result.CodeResults) < RESULTS_PER_PAGE {
		return
	}

	// This rate limit is very conservative, it abides by the current
	// unauthenticated limit imposed by GitHub (10 requests per minute).
	// In theory, we could lower this to 2 since we authenticate.
	time.Sleep(6 * time.Second)

	fetchResults(client, ctx, userName, page + 1)
}

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

	fetchResults(client, ctx, userName,1)
}
