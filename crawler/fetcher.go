package crawler

import "fmt"

type Fetcher interface {
	// Fetch returns the body of URL and a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type GolangMockFetcher map[string]*FetchResult

func NewGolangMockFetcher() GolangMockFetcher {
	return GolangMockFetcher{
		"https://golang.org/": NewFetchResult(
			"The Go Programming Language",
			[]string{
				"https://golang.org/pkg/",
				"https://golang.org/cmd/",
			}),
		"https://golang.org/pkg/": NewFetchResult(
			"Packages",
			[]string{
				"https://golang.org/",
				"https://golang.org/cmd/",
				"https://golang.org/pkg/fmt/",
				"https://golang.org/pkg/os/",
			}),
		"https://golang.org/pkg/fmt/": NewFetchResult(
			"Package fmt",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			}),
		"https://golang.org/pkg/os/": NewFetchResult(
			"Package os",
			[]string{
				"https://golang.org/",
				"https://golang.org/pkg/",
			}),
	}
}

func (f GolangMockFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.Body, res.URLs, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}
