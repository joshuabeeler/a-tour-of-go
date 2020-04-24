package crawler

type FetchResult struct {
	Body string
	URLs []string
}

func NewFetchResult(body string, urls []string) *FetchResult {
	return &FetchResult{
		Body: body,
		URLs: urls,
	}
}
