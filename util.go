package main

func contains(haystack []SearchResult, url string) bool {
	for _, candidate := range haystack {
		if url == candidate.URL {
			return true
		}
	}
	return false
}
