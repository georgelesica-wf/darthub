package main

import (
	"github.com/google/go-github/github"
	"time"
	"context"
	"fmt"
)

const RESULTS_PER_PAGE = 100

type GithubSearcher struct {
	client *github.Client
}

func (g *GithubSearcher) Search(params SearchParams) ([]SearchResult, error) {
	ctx := context.Background()

	return g.search(&params, ctx, make([]SearchResult, 0))
}

func (g *GithubSearcher) search(params *SearchParams, ctx context.Context, results []SearchResult) ([]SearchResult, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page: params.Page,
			PerPage: RESULTS_PER_PAGE,
		},
	}

	query := fmt.Sprintf("user:%s filename:%s extension:%s",
		params.User, params.Filename, params.Extension)

	githubResult, _, err := g.client.Search.Code(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	for _, item := range githubResult.CodeResults {
		repoURL := item.Repository.GetHTMLURL()

		if contains(results, repoURL) {
			continue
		}

		results = append(results, SearchResult{URL: repoURL})
	}

	// We are finished when we get back fewer than the max results.
	if len(githubResult.CodeResults) < RESULTS_PER_PAGE {
		return results, nil
	}

	// This rate limit is very conservative, it abides by the current
	// unauthenticated limit imposed by GitHub (10 requests per minute).
	// In theory, we could lower this to 2 since we authenticate.
	time.Sleep(6 * time.Second)

	params.Page++
	return g.search(params, ctx, results)
}
