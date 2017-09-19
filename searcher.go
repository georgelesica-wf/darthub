package main

type SearchParams struct {
	Extension string
	Filename  string
	Page      int
	PerPage   int
	User      string
}

type SearchResult struct {
	URL        string
	TrimmedURL string
}

type Searcher interface {
	Search(params *SearchParams) ([]SearchResult, error)
}
