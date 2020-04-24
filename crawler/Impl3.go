package crawler

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// JOSH'S CODE
////////////////////////////////////////////////////////////////////////////////////////////////////

type Impl3 struct {
	fetcher Fetcher				// We use this to make "web" fetches.
	chWorkers chan workerResult	// The broker receives new URLs from workers via this channel.
	inFlight int				// Number of workers doing work.
	seenMux sync.Mutex			// Used to synchronize access to the seen map.
	seen map[string]bool		// Collection of URLs we've already crawled.
}

// Launch a parallel crawl operation.
func (c *Impl3) Begin(crawlerDone chan bool, url string, depth int, fetcher Fetcher) {
	c.fetcher = fetcher
	c.Seen(url)
	c.Launch(url, depth)
	go c.Broker(crawlerDone)
}

// Goroutine that manages all of the workers involved in a parallel crawl operation.
func (c *Impl3) Broker(crawlerDone chan bool) {
	//fmt.Println("broker launch")
	
	for ; c.inFlight > 0; {
		// Blocking is OK. If we have any workers in flight, we're guaranteed to get results.
		res := <- c.chWorkers
		
		// This is not a launchable fetch.
		if res.depth > 0 {
			c.Launch(res.url, res.depth)
		}
		
		// Worker is finished.
		if res.done {
			c.inFlight = c.inFlight - 1
		}
	}
	
	//fmt.Println("broker finished")
	crawlerDone <- true
}

// Check if a URL has been seen before (and add it to the seen map if it hasn't been) and return the result.
func (c *Impl3) Seen(url string) bool {
	c.seenMux.Lock()
	_, ok := c.seen[url]
	if !ok {
		c.seen[url] = true
	}
	c.seenMux.Unlock()
	return ok
}

// Start a new sub-crawl at a given URL.
func (c *Impl3) Launch(url string, depth int) {
	c.inFlight = c.inFlight + 1
	//fmt.Printf("worker launch: %s\n", req.url)
	go c.Crawl(url, depth, c.fetcher)
}

// Fetch a URL and issue sub-crawl requests for all of the URLs children.
func (c *Impl3) Crawl(url string, depth int, fetcher Fetcher) {
	//fmt.Println("crawling")
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		c.chWorkers <- workerResult{ "", 0, true }
		return
	}
	
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !c.Seen(u) {
			//fmt.Printf("adding: %s\n", u)
			c.chWorkers <- workerResult{ u, depth - 1, false }
		}
	}
	
	c.chWorkers <- workerResult{ "", 0, true }
}

func main() {
	//Crawl("https://golang.org/", 4, fetcher)

	crawler := Impl3{
		seen: make(map[string]bool),
		chWorkers: make(chan workerResult, 128),
	}

	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, fetcher)
	<- done
}
