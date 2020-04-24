package crawler

import "testing"

func TestImpl1(t *testing.T) {
	c := Impl1{
		seen: make(map[string]bool),
		chReq: make(chan request, 128),
		chDone: make(chan int, 128),	// why is , needed here?
	}

	done := make(chan bool)
	c.Begin(done, "https://golang.org/", 4, NewGolangMockFetcher())
	<- done
}

func TestImpl2(t *testing.T) {
	c := Impl2{
		seen: make(map[string]bool),
		chWorkers: make(chan workerResult, 128),	// why is , needed here?
	}

	done := make(chan bool)
	c.Begin(done, "https://golang.org/", 4, NewGolangMockFetcher())
	<- done
}

func TestImpl3(t *testing.T) {
	crawler := NewImpl3()

	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, NewGolangMockFetcher())
	<- done

	expectedResults := []string{
		"https://golang.org/",
		"https://golang.org/cmd/",
		"https://golang.org/pkg/",
		"https://golang.org/pkg/fmt/",
		"https://golang.org/pkg/os/",
	}

	actualResults := crawler.GetResults()

	for _, er := range expectedResults {
		found := false
		for _, ar := range actualResults {
			if er == ar {
				found = true;
				break;
			}
		}
		if !found {
			t.Errorf("Expected to find %s in results, but didn't.", er)
		}
	}
}
