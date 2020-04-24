package crawler

import "testing"

func TestImpl3(t *testing.T) {
	crawler := GetImpl3()

	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, fetcher)
	<- done
}
