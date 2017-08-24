package main

import (
	"github.com/google/go-github/github"
	"context"
	"golang.org/x/oauth2"
	"os"
)

var repoURLs []string = make([]string, 0)

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if needle == candidate {
			return true
		}
	}
	return false
}

func fetchResults(client *github.Client, ctx context.Context, page int) {
	opts := &github.SearchOptions{ListOptions: github.ListOptions{Page: page, PerPage: 100}}

	// list all repositories for the authenticated user
	result, _, err := client.Search.Code(ctx, "user:Workiva filename:pubspec extension:yaml", opts)
	if err != nil {
		panic(err)
	}

	if len(result.CodeResults) == 0 {
		return
	}

	for _, item := range result.CodeResults {
		repoURL := item.Repository.GetHTMLURL()

		if contains(repoURLs, repoURL) {
			continue
		}

		repoURLs = append(repoURLs, repoURL)
		println(repoURL)
	}

	fetchResults(client, ctx, page + 1)
}

func main() {
	token, present := os.LookupEnv("DARTHUB_TOKEN")
	if !present {
		panic("DARTHUB_TOKEN is not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	fetchResults(client, ctx, 1)
}
