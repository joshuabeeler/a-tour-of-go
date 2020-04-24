package main

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// JOSH'S CODE
////////////////////////////////////////////////////////////////////////////////////////////////////

type Crawler struct {
	fetcher Fetcher				// We use this to make "web" fetches.
	chWorkers chan WorkerResult	// The broker receives new URLs from workers via this channel.
	inFlight int				// Number of workers doing work.
	seenMux sync.Mutex			// Used to synchronize access to the seen map.
	seen map[string]bool		// Collection of URLs we've already crawled.
}

type WorkerResult struct {
	url string			// URL to fetch. "" if nothing to fetch.
	depth int			// Stop fetching when this reaches zero.
	done bool			// True when a worker is done, false otherwise.
}

// Launch a parallel crawl operation.
func (c *Crawler) Begin(crawlerDone chan bool, url string, depth int, fetcher Fetcher) {
	c.fetcher = fetcher
	c.Seen(url)
	c.Launch(url, depth)
	go c.Broker(crawlerDone)
}

// Goroutine that manages all of the workers involved in a parallel crawl operation.
func (c *Crawler) Broker(crawlerDone chan bool) {
	//fmt.Println("broker launch")
	
	for ; c.inFlight > 0; {
		// Blocking is OK. If we have any workers in flight, we're guaranteed to get results.
		res := <- c.chWorkers
		
		// This is not a launchable fetch.
		if (res.depth > 0) {
			c.Launch(res.url, res.depth)
		}
		
		// Worker is finished.
		if (res.done) {
			c.inFlight = c.inFlight - 1
		}
	}
	
	//fmt.Println("broker finished")
	crawlerDone <- true
}

// Check if a URL has been seen before (and add it to the seen map if it hasn't been) and return the result.
func (c *Crawler) Seen(url string) bool {
	c.seenMux.Lock()
	_, ok := c.seen[url]
	if !ok {
		c.seen[url] = true
	}
	c.seenMux.Unlock()
	return ok
}

// Start a new sub-crawl at a given URL.
func (c *Crawler) Launch(url string, depth int) {
	c.inFlight = c.inFlight + 1
	//fmt.Printf("worker launch: %s\n", req.url)
	go c.Crawl(url, depth, c.fetcher)
}

// Fetch a URL and issue sub-crawl requests for all of the URLs children.
func (c *Crawler) Crawl(url string, depth int, fetcher Fetcher) {
	//fmt.Println("crawling")
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		c.chWorkers <- WorkerResult{ "", 0, true }
		return
	}
	
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !c.Seen(u) {
			//fmt.Printf("adding: %s\n", u)
			c.chWorkers <- WorkerResult{ u, depth - 1, false }
		}
	}
	
	c.chWorkers <- WorkerResult{ "", 0, true }
}

var crawler = Crawler{
	seen: make(map[string]bool),
	chWorkers: make(chan WorkerResult, 128),	// why is , needed here?
}

func main() {
	//Crawl("https://golang.org/", 4, fetcher)
	
	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, fetcher)
	<- done
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// EXERCISE FRAMEWORK CODE
////////////////////////////////////////////////////////////////////////////////////////////////////

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	
	if depth <= 0 {
		return
	}
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}
	return
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
