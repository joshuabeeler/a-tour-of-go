package crawler

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// JOSH'S CODE
////////////////////////////////////////////////////////////////////////////////////////////////////

type Crawler1 struct {
	fetcher Fetcher			// We use this to make "web" fetches.
	chReq chan Request		// The broker receives new URLs from workers via this channel.
	chDone chan int			// The broker gets notified when workers are finished via this channel.
	inFlight int			// Number of workers doing work.
	seenMux sync.Mutex		// Used to synchronize access to the seen map.
	seen map[string]bool	// Collection of URLs we've already crawled.
}

type Request struct {
	url string			// URL to fetch.
	depth int			// Stop fetching when this reaches zero.
}

// Launch a parallel crawl operation.
func (c *Crawler1) Begin(crawlerDone chan bool, url string, depth int, fetcher Fetcher) {
	c.fetcher = fetcher
	c.Seen(url)
	c.Launch(Request{url, depth})
	go c.Broker(crawlerDone)
}

// Goroutine that manages all of the workers involved in a parallel crawl operation.
func (c *Crawler1) Broker(crawlerDone chan bool) {
	//fmt.Println("broker launch")
	
	for ; c.inFlight > 0; {
		select {
			case <- c.chDone:
				//fmt.Println("worker finished")
				c.inFlight = c.inFlight - 1
			case req := <- c.chReq:
				c.Launch(req)
		}
	}
	
	//fmt.Println("broker finished")
	crawlerDone <- true
}

// Check if a URL has been seen before (and add it to the seen map if it hasn't been) and return the result.
func (c *Crawler1) Seen(url string) bool {
	c.seenMux.Lock()
	_, ok := c.seen[url]
	if !ok {
		c.seen[url] = true
	}
	c.seenMux.Unlock()
	return ok
}

// Start a new sub-crawl at a given URL.
func (c *Crawler1) Launch(req Request) {
	c.inFlight = c.inFlight + 1
	//fmt.Printf("worker launch: %s\n", req.url)
	go c.Crawl(req.url, req.depth, c.fetcher)
}

// Fetch a URL and issue sub-crawl requests for all of the URLs children.
func (c *Crawler1) Crawl(url string, depth int, fetcher Fetcher) {
	//fmt.Println("crawling")
	
	if depth <= 0 {
		c.chDone <- 1
		return
	}
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		c.chDone <- 1
		return
	}
	
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !c.Seen(u) {
			//fmt.Printf("adding: %s\n", u)
			c.chReq <- Request{ u, depth - 1 }
		}
	}
	
	c.chDone <- 1
}

func main1() {
	//Crawl("https://golang.org/", 4, fetcher)

	crawler := Crawler1{
		seen: make(map[string]bool),
		chReq: make(chan Request, 128),
		chDone: make(chan int, 128),	// why is , needed here?
	}
	
	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, fetcher)
	<- done
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// EXERCISE FRAMEWORK CODE
////////////////////////////////////////////////////////////////////////////////////////////////////

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
/*
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
*/
