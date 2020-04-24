package main

import "fmt"
import "github.com/joshuabeeler/a-tour-of-go/crawler"

func main() {
	fmt.Println("Hello, Josh.")

	c := crawler.NewImpl3()

	done := make(chan bool)
	c.Begin(done, "https://golang.org/", 4, crawler.NewGolangMockFetcher())
	<- done

	fmt.Println("Bye, Josh.")
}
