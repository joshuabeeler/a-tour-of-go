package crawler

import "testing"

func TestImpl3(t *testing.T) {
	crawler := GetImpl3()

	done := make(chan bool)
	crawler.Begin(done, "https://golang.org/", 4, fetcher)
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
